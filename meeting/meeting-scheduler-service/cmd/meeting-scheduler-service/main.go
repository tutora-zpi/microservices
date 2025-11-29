package main

import (
	"context"
	"log"
	"meeting-scheduler-service/internal/app/interfaces"
	"meeting-scheduler-service/internal/app/usecase"
	"meeting-scheduler-service/internal/config"
	"meeting-scheduler-service/internal/infrastructure/bus"
	"meeting-scheduler-service/internal/infrastructure/messaging"
	"meeting-scheduler-service/internal/infrastructure/postgres"
	"meeting-scheduler-service/internal/infrastructure/redis"
	"meeting-scheduler-service/internal/infrastructure/repository"
	"meeting-scheduler-service/internal/infrastructure/rest"
	"meeting-scheduler-service/internal/infrastructure/rest/v1/handlers"
	"meeting-scheduler-service/internal/infrastructure/security"
	"meeting-scheduler-service/internal/infrastructure/server"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	redisdb "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	security.FetchSignKey()
}

func main() {
	var redisClient *redisdb.Client
	var postgresClient *gorm.DB
	var postgresConfig postgres.PostgresConfig = *postgres.NewPostgresConfig(time.Second * 5)
	var redisConfig redis.RedisConfig = *redis.NewRedisConfig(time.Second * 5)
	var rabbitmqConfig messaging.RabbitConfig = *messaging.NewRabbitMQConfig(time.Second*5, 10)
	var broker interfaces.Broker
	var wg sync.WaitGroup

	var errors chan error = make(chan error, 2)

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dispatcher := bus.NewDispatcher()
	wg.Go(func() {
		var err error
		broker, err = messaging.NewRabbitBroker(rabbitmqConfig, dispatcher)
		if err != nil {
			errors <- err
		}
	})

	wg.Go(func() {
		var err error
		redisClient, err = redis.NewRedis(rootCtx, redisConfig)
		if err != nil {
			errors <- err
		}

	})

	wg.Go(func() {
		var err error
		postgresClient, err = postgres.NewPostgres(rootCtx, postgresConfig)
		if err != nil {
			errors <- err
		}

	})

	wg.Wait()
	close(errors)

	for err := range errors {
		log.Fatalf("Error during services init: %v", err)
	}

	closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer func() {
		postgres.Close(closeCtx, postgresClient, postgresConfig)
		redis.Close(closeCtx, redisClient, redisConfig)
		broker.Close()

		cancel()
	}()

	meetingRepo := repository.NewMeetingRepository(redisClient)
	plannedMeetingsRepo := repository.NewPlannedMeetingsRepository(postgresClient)

	terminator := usecase.NewMeetingTerminator()

	meetingManager := usecase.NewManageMeeting(
		broker,
		meetingRepo,
		plannedMeetingsRepo,
		terminator,
		rabbitmqConfig.MeetingExchange,
	)

	planner := usecase.NewPlanner(rootCtx, meetingManager, usecase.PlannerConfig{
		FetchIntervalMinutes: 1,
	})

	go terminator.Run(rootCtx, meetingManager.Stop)
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
