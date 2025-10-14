package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	generalHandler "signaling-service/internal/app/socket_event_handler/general"
	generalDomain "signaling-service/internal/domain/ws_event/general"
	"signaling-service/internal/infrastructure/bus"
	"signaling-service/internal/infrastructure/cache"
	"signaling-service/internal/infrastructure/rest"
	"signaling-service/internal/infrastructure/rest/v1/handlers"
	"signaling-service/internal/infrastructure/server"
	"signaling-service/internal/infrastructure/ws"
	"syscall"
)

func init() {

}

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// initCtx, cancel := context.WithTimeout(rootCtx, 3*time.Second)
	// defer cancel()

	hub := ws.NewHub()

	dispacher := bus.NewDispatcher()

	dispacher.Register(&generalDomain.UserJoinedEvent{}, generalHandler.NewUserJoinedHandler(hub))
	dispacher.Register(&generalDomain.UserLeftEvent{}, generalHandler.NewUserLeftHandler(hub))
	dispacher.Register(&generalDomain.HeartbeatEvent{}, generalHandler.NewHeartbeatHandler(hub))

	tokenService := cache.NewTokenService(nil)

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
