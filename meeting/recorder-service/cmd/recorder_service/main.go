package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	eventhandlers "recorder-service/internal/app/event_handlers"
	factoryimpl "recorder-service/internal/app/factory_impl"
	"recorder-service/internal/app/service"
	"recorder-service/internal/config"
	"recorder-service/internal/domain/event"
	"recorder-service/internal/domain/ws_event/rtc"
	"recorder-service/internal/infrastructure/bus"
	"recorder-service/internal/infrastructure/cache"
	"recorder-service/internal/infrastructure/mongo"
	"recorder-service/internal/infrastructure/redis"
	repoimpl "recorder-service/internal/infrastructure/repository"
	"recorder-service/internal/infrastructure/rest"
	"recorder-service/internal/infrastructure/rest/v1/handlers"
	"recorder-service/internal/infrastructure/security"
	"recorder-service/internal/infrastructure/server"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	redisdb "github.com/redis/go-redis/v9"
	mongodb "go.mongodb.org/mongo-driver/v2/mongo"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	security.FetchSignKey()
}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var errorCh chan error

	var mongoClient *mongodb.Client
	var redisClient *redisdb.Client

	var mongoConfig mongo.MongoConfig = *mongo.NewMongoConfig(time.Second * 5)
	var redisConfig redis.RedisConfig = *redis.NewRedisConfig(time.Second * 5)

	dispatcher := bus.NewDispatcher()

	var wg sync.WaitGroup

	wg.Go(func() {
		var err error
		mongoClient, err = mongo.NewMongoClient(rootCtx, mongoConfig)

		errorCh <- err
	})

	wg.Go(func() {
		var err error
		redisClient, err = redis.NewRedis(rootCtx, redisConfig)

		errorCh <- err
	})

	wg.Wait()

	close(errorCh)

	for err := range errorCh {
		log.Fatalf("Something went wrong with services: %v", err)
	}

	voiceRepo := repoimpl.NewVoiceMeetingRepository(mongoClient, mongoConfig)
	botRepo := cache.NewBotRepository(redisClient)

	voiceSessionService := service.NewVoiceSessionService(voiceRepo)

	clientFactory := factoryimpl.NewClientFactory()
	recorderFactory := factoryimpl.NewRecorderFactory()

	botService := service.NewBotService(botRepo, recorderFactory, clientFactory)

	dispatcher.Register(&event.MeetingStartedEvent{}, eventhandlers.NewMeetingStartedHandler(voiceRepo))
	dispatcher.Register(&event.StopRecordingMeetingEvent{}, eventhandlers.NewStopRecordingMeetingHandler(botService, voiceRepo))
	dispatcher.Register(&event.RecordMeetingEvent{}, eventhandlers.NewRecorderMeetingHandler(botService))

	dispatcher.Register(&rtc.AnswerWSEvent{}, eventhandlers.NewAnswerHandler(botService))
	dispatcher.Register(&rtc.IceCandidateWSEvent{}, eventhandlers.NewIceCandidateHandler(botService))

	handlers := handlers.NewHandler(voiceSessionService)

	router := rest.NewRouter(handlers)
	server := server.NewServer(router)
	go func() {
		if err := server.StartAndListen(); err != nil {
			log.Printf("server stopped with error: %v", err)
			stop()
		}
	}()

	<-rootCtx.Done()

	if err := server.GracefulShutdown(context.Background()); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
