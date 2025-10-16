package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	boardHandler "signaling-service/internal/app/socket_event_handler/board"
	chatHandler "signaling-service/internal/app/socket_event_handler/chat"
	generalHandler "signaling-service/internal/app/socket_event_handler/general"
	rtcHandler "signaling-service/internal/app/socket_event_handler/rtc"
	"signaling-service/internal/config"
	boardDomain "signaling-service/internal/domain/ws_event/board"
	chatDomain "signaling-service/internal/domain/ws_event/chat"
	generalDomain "signaling-service/internal/domain/ws_event/general"
	rtcDomain "signaling-service/internal/domain/ws_event/rtc"
	"signaling-service/internal/infrastructure/bus"
	"signaling-service/internal/infrastructure/cache/repo"
	"signaling-service/internal/infrastructure/cache/service"
	"signaling-service/internal/infrastructure/messaging"
	"signaling-service/internal/infrastructure/redis"
	"signaling-service/internal/infrastructure/rest"
	"signaling-service/internal/infrastructure/rest/v1/handlers"
	securityRepo "signaling-service/internal/infrastructure/security/repository"
	"signaling-service/internal/infrastructure/server"
	"signaling-service/internal/infrastructure/ws"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var rabbitmqConfig messaging.RabbitConfig
var redisConfig redis.RedisConfig

func setupRedisConfig() {
	db := 0
	if v := os.Getenv(config.REDIS_DB); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			db = parsed
		}
	}

	redisConfig = redis.RedisConfig{
		Addr:     os.Getenv(config.REDIS_ADDR),
		Password: os.Getenv(config.REDIS_PASSWORD),
		DB:       db,
	}
}

func setupRabbitMQConfig() {
	rabbitmqConfig = messaging.RabbitConfig{
		User:    os.Getenv(config.RABBITMQ_DEFAULT_USER),
		Pass:    os.Getenv(config.RABBITMQ_DEFAULT_PASS),
		Host:    os.Getenv(config.RABBITMQ_HOST),
		Port:    os.Getenv(config.RABBITMQ_PORT),
		URL:     os.Getenv(config.RABBITMQ_URL),
		Retries: 3,

		ExchangeNames: []string{
			os.Getenv(config.CHAT_EXCHANGE),
			os.Getenv(config.BOARD_EXCHANGE),
		},
	}
}

func init() {

	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		if err := godotenv.Load(".env.local"); err != nil {
			log.Panic(".env* file not found. Please check path or provide one.")
		}
	}

	setupRabbitMQConfig()
	setupRedisConfig()

	// securityJWT.FetchSignKey()
}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initCtx, cancel := context.WithTimeout(rootCtx, 3*time.Second)
	defer cancel()

	redis, err := redis.NewRedis(initCtx, redisConfig)
	if err != nil {
		log.Panicf("Failed to create redis: %v", err)
	}
	defer func() {
		log.Println("Closing redis connection...")
		redis.Close()
	}()

	hub := ws.NewHub()
	cacheRepo := repo.NewCacheEventRepository(redis, 10, time.Minute)
	cacheService := service.NewCacheEventSerivce(cacheRepo)

	dispacher := bus.NewDispatcher()
	rabbitmq, err := messaging.NewRabbitBroker(rabbitmqConfig, dispacher)
	if err != nil {
		log.Panicf("Failed to create rabbitmq: %v", err)
	}
	defer rabbitmq.Close()

	eventBuffer := bus.NewEventBuffer(rabbitmq)

	go eventBuffer.Work(rootCtx)

	dispacher.Register(&generalDomain.UserJoinedEvent{}, generalHandler.NewUserJoinedHandler(hub, cacheService))
	dispacher.Register(&generalDomain.UserLeftEvent{}, generalHandler.NewUserLeftHandler(hub))
	dispacher.Register(&chatDomain.UserTypingEvent{}, chatHandler.NewUserTypingHandler(hub))
	dispacher.Register(&chatDomain.SendMessageEvent{}, chatHandler.NewSendMessageHandler(hub, eventBuffer, cacheService))
	dispacher.Register(&boardDomain.BoardUpdateEvent{}, boardHandler.NewBoardUpdateHandler(hub, eventBuffer, cacheService))
	dispacher.Register(&rtcDomain.AnswerEvent{}, rtcHandler.NewAnswerHandler(hub))
	dispacher.Register(&rtcDomain.IceCandidateEvent{}, rtcHandler.NewIceCandidateHandler(hub))
	dispacher.Register(&rtcDomain.OfferEvent{}, rtcHandler.NewOfferHandler(hub))

	tokenService := securityRepo.NewTokenService(redis)

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
