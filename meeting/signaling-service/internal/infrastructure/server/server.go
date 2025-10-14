package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"signaling-service/internal/config"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Server struct {
	s    *http.Server
	host string
}

const DEFAULT_PORT string = "8010"
const DEFAULT_HOST string = "localhost"

func NewServer(router *mux.Router) *Server {
	port := os.Getenv(config.APP_PORT)

	host := DEFAULT_HOST

	if port == "" {
		port = DEFAULT_PORT
	}

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"DELETE", "GET", "HEAD", "OPTIONS", "POST"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := corsHandler.Handler(router)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{s: s, host: host}
}

func (apiServer *Server) GracefulShutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	log.Println("Shutting down gracefully...")
	if err := apiServer.s.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited")
	return nil
}

func (apiServer *Server) StartAndListen() error {

	log.Printf("Server is listening on: http://%s%s", apiServer.host, apiServer.s.Addr)

	if err := apiServer.s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server error: %w", err)
	}
	return nil
}
