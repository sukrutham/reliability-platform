package main

import (
	"log"
	"net/http"

	"github.com/sukrutham/reliability-platform/internal/handlers"
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

	addr := ":8080"
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
