package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"notification-serivce/internal/config"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	s *http.Server
}

const DEFAULT_PORT string = "8888"
const DEFAULT_HOST string = "localhost"

func NewServer(router *mux.Router) *Server {
	port := os.Getenv(config.APP_PORT)
	host := os.Getenv(config.APP_ENV)

	if host == "" {
		host = DEFAULT_HOST
	}

	if port == "" {
		port = DEFAULT_PORT
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "OPTIONS", "HEAD"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := corsHandler.Handler(router)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{s: s}
}

// GracefulShutdown listens for system signals and shuts down the server cleanly,
// allowing up to 5 seconds for open connections to finish.
func (apiServer *Server) GracefulShutdown(done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.s.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func (apiServer *Server) StartAndListen() {
	log.Printf("Server is listening on: http://%s%s", os.Getenv(config.APP_ENV), apiServer.s.Addr)
	err := apiServer.s.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Panicf("HTTP server error: %s", err)
	}
}
