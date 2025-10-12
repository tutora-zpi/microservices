package main

import (
	"context"
	"log"
	"meeting-scheduler-service/internal/app/usecase"
	"meeting-scheduler-service/internal/config"
	"meeting-scheduler-service/internal/infrastructure/bus"
	"meeting-scheduler-service/internal/infrastructure/messaging"
	"meeting-scheduler-service/internal/infrastructure/postgres"
	"meeting-scheduler-service/internal/infrastructure/redis"
	"meeting-scheduler-service/internal/infrastructure/rest"
	"meeting-scheduler-service/internal/infrastructure/rest/v1/handlers"
	"meeting-scheduler-service/internal/infrastructure/security"
	"meeting-scheduler-service/internal/infrastructure/server"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

var postgresConfig postgres.PostgresConfig
var redisConfig redis.RedisConfig
var rabbitmqConfig messaging.RabbitConfig

func setupPostgresConfig() {
	sslMode, err := strconv.ParseBool(os.Getenv(config.POSTGRES_SSLMODE))
	if err != nil {
		sslMode = false
	}

	postgresConfig = postgres.PostgresConfig{
		User:     os.Getenv(config.POSTGRES_USER),
		Password: os.Getenv(config.POSTGRES_PASS),

		Host: os.Getenv(config.POSTGRES_HOST),
		Port: os.Getenv(config.POSTGRES_PORT),

		DatabaseName: os.Getenv(config.POSTGRES_DBNAME),
		SSLMode:      sslMode,
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
			os.Getenv(config.NOTIFICATION_EXCHANGE_QUEUE_NAME),
			os.Getenv(config.EVENT_EXCHANGE_QUEUE_NAME),
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

	setupPostgresConfig()
	setupRabbitMQConfig()

	security.FetchSignKey()
}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	initCtx, cancel := context.WithTimeout(rootCtx, 3*time.Second)
	defer cancel()

	// BROKER
	dispatcher := bus.NewDispatcher()
	broker, err := messaging.NewRabbitBroker(rabbitmqConfig, dispatcher)
	if err != nil {
		log.Panicf("Failed to create broker: %v", err)
	}

	defer broker.Close()

	// REDIS REPO
	meetingRepo, err := redis.NewMeetingRepo(initCtx, redisConfig)
	if err != nil {
		log.Panicf("Failed to create redis repo: %v", err)
	}
	defer meetingRepo.Close()

	plannedMeetingsRepo, err := postgres.NewMeetingsRepository(postgresConfig)
	if err != nil {
		log.Panicf("Failed to create postgres repo: %v", err)
	}
	defer plannedMeetingsRepo.Close()

	// APP
	meetingManager := usecase.NewManageMeeting(broker, meetingRepo, plannedMeetingsRepo, os.Getenv(config.NOTIFICATION_EXCHANGE_QUEUE_NAME), os.Getenv(config.EVENT_EXCHANGE_QUEUE_NAME))

	// PLANNER
	planner := usecase.NewPlanner(meetingManager, usecase.PlannerConfig{
		FetchIntervalMinutes: 1,
	})

	go planner.Listen(rootCtx)
	go planner.RerunNotStartedMeetings(rootCtx)

	// HTTP
	router := rest.NewRouter(handlers.NewManageMeetingHandler(meetingManager))
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
