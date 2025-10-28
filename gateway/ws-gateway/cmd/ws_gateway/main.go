package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	eventhandler "ws-gateway/internal/app/event_handler"
	"ws-gateway/internal/app/interfaces"
	boardHandler "ws-gateway/internal/app/ws_event_handler/board"
	chatHandler "ws-gateway/internal/app/ws_event_handler/chat"
	generalHandler "ws-gateway/internal/app/ws_event_handler/general"
	rtcHandler "ws-gateway/internal/app/ws_event_handler/rtc"
	"ws-gateway/internal/config"
	"ws-gateway/internal/domain/event"
	boardDomain "ws-gateway/internal/domain/ws_event/board"
	chatDomain "ws-gateway/internal/domain/ws_event/chat"
	generalDomain "ws-gateway/internal/domain/ws_event/general"
	rtcDomain "ws-gateway/internal/domain/ws_event/rtc"
	"ws-gateway/internal/infrastructure/bus"
	"ws-gateway/internal/infrastructure/cache/repo"
	"ws-gateway/internal/infrastructure/cache/service"
	"ws-gateway/internal/infrastructure/messaging"
	redisdb "ws-gateway/internal/infrastructure/redis"
	"ws-gateway/internal/infrastructure/rest"
	"ws-gateway/internal/infrastructure/rest/v1/handlers"
	security "ws-gateway/internal/infrastructure/security/jwt"
	securityRepo "ws-gateway/internal/infrastructure/security/repository"
	"ws-gateway/internal/infrastructure/server"
	"ws-gateway/internal/infrastructure/ws"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func init() {

	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	security.FetchSignKey()
}

func main() {
	var wg sync.WaitGroup
	var broker interfaces.Broker
	var redisClient *redis.Client
	var rabbitmqConfig messaging.RabbitConfig = *messaging.NewRabbitMQConfig(time.Second*5, 10)
	var redisConfig redisdb.RedisConfig = *redisdb.NewRedisConfig(time.Second * 5)
	var errors chan error = make(chan error, 2)
	dispacher := bus.NewDispatcher()

	hub := ws.NewHub()
	defer hub.Close()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initCtx, cancel := context.WithTimeout(rootCtx, 3*time.Second)
	defer cancel()

	wg.Go(func() {
		var err error
		redisClient, err = redisdb.NewRedis(initCtx, redisConfig)
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		broker, err = messaging.NewRabbitBroker(rabbitmqConfig, dispacher)
		if err != nil {
			errors <- err
		}
	})

	wg.Wait()
	close(errors)

	for err := range errors {
		log.Fatalf("Error: %v", err)
	}

	defer func() {
		if redisClient != nil {
			log.Println("Closing redis connection...")
			_ = redisClient.Close()
		}
		if broker != nil {
			broker.Close()
		}
	}()

	cacheRepo := repo.NewCacheEventRepository(redisClient, 10, time.Second*30)
	cacheService := service.NewCacheEventSerivce(cacheRepo)

	eventBuffer := bus.NewEventBuffer(broker)
	defer eventBuffer.Close()

	dispacher.Register(&generalDomain.UserJoinedWSEvent{}, generalHandler.NewUserJoinedHandler(hub, cacheService))
	dispacher.Register(&generalDomain.UserLeftWSEvent{}, generalHandler.NewUserLeftHandler(hub))
	dispacher.Register(&chatDomain.UserTypingWSEvent{}, chatHandler.NewUserTypingHandler(hub))
	dispacher.Register(&chatDomain.SendMessageWSEvent{}, chatHandler.NewSendMessageHandler(hub, eventBuffer, cacheService))
	dispacher.Register(&boardDomain.BoardUpdateWSEvent{}, boardHandler.NewBoardUpdateHandler(hub, eventBuffer, cacheService))
	dispacher.Register(&rtcDomain.AnswerWSEvent{}, rtcHandler.NewAnswerHandler(hub))
	dispacher.Register(&rtcDomain.IceCandidateWSEvent{}, rtcHandler.NewIceCandidateHandler(hub))
	dispacher.Register(&rtcDomain.OfferWSEvent{}, rtcHandler.NewOfferHandler(hub))
	dispacher.Register(&event.SendFileMessageEvent{}, eventhandler.NewSendFileMessageHandler(hub, cacheService))

	go eventBuffer.Work(rootCtx)
	go func() {
		err := broker.Consume(rootCtx, rabbitmqConfig.FileQueue)
		if err != nil {
			log.Printf("Consuming error: %v", err)
		}
	}()

	tokenService := securityRepo.NewTokenService(redisClient)

	handlers := handlers.NewHandlers(dispacher, hub, tokenService)

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
