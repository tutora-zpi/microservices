package main

import (
	"log"
	handlers "notification-serivce/internal/infrastructure/rest"
	"notification-serivce/internal/infrastructure/security"
	"notification-serivce/internal/infrastructure/server"
	"notification-serivce/internal/infrastructure/sse"

	"github.com/joho/godotenv"
)

func init() {

	// loading envs
	if err := godotenv.Load(); err != nil {
		log.Panic(".env* file not found. Please check path or provide one.")
	}

	// fetch jwt secret
	security.FetchSignKey()
}

func main() {
	sseManager := sse.NewSSEManager()

	server := server.NewServer(handlers.NewRouter(sseManager))

	go func() {
		server.StartAndListen()
	}()

	done := make(chan bool, 1)

	go func() {
		server.GracefulShutdown(done)
	}()

	<-done
}
