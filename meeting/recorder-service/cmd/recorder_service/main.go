package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"recorder-service/internal/config"
	"recorder-service/internal/infrastructure/rest"
	"recorder-service/internal/infrastructure/server"
	"syscall"

	"github.com/joho/godotenv"
)

func init() {
	env := os.Getenv(config.APP_ENV)

	if env == "" || env == "localhost" || env == "127.0.0.1" {
		_ = godotenv.Load(".env.local")
	}

	// security.FetchSignKey()
}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	router := rest.NewRouter()
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
