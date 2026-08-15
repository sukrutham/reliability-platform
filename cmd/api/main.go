package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sukrutham/reliability-platform/internal/handlers"
	"github.com/sukrutham/reliability-platform/internal/metrics"
	"github.com/sukrutham/reliability-platform/internal/store"
)

func main() {
	todoStore := store.NewTodoStore()
	todoHandler := handlers.NewTodoHandler(todoStore)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Healthz)
	mux.HandleFunc("GET /readyz", handlers.Readyz)
	mux.HandleFunc("GET /api/v1/todos", todoHandler.List)
	mux.HandleFunc("POST /api/v1/todos", todoHandler.Create)
	mux.HandleFunc("GET /api/v1/todos/{id}", todoHandler.Get)
	mux.Handle("GET /metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: metrics.Middleware(mux),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	stop()
	log.Println("shutdown signal received, draining in-flight requests")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("graceful shutdown failed: %v", err)
	}
	log.Println("server stopped cleanly")
}
