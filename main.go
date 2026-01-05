package main

import (
	"assignment2/internal/api"
	"assignment2/internal/model"
	"assignment2/internal/queue"
	"assignment2/internal/store"
	"assignment2/internal/worker"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("starting server...")

	// initialization
	repository := store.NewRepository[string, *model.Task]()
	taskQueue := queue.NewTaskQueue(100)
	workerPool := worker.NewPool(taskQueue, repository)
	handler := api.NewHandler(repository, taskQueue)

	// worker pool with 2 workers
	workerPool.Start(2)

	// setup HTTP routes
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/tasks" {
			if r.Method == http.MethodPost {
				handler.CreateTask(w, r)
			} else if r.Method == http.MethodGet {
				handler.GetTasks(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			handler.GetTask(w, r)
		}
	})

	http.HandleFunc("/stats", handler.GetStats)

	// create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}

	// start server in goroutine
	go func() {
		log.Println("Server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// shutdown mesage
	<-sigCh
	log.Println("\n starting graceful shutdown...")

	// stop accepting new requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
	log.Println("HTTP server stopped")

	// stopping workers
	workerPool.Stop()

	log.Println("graceful shutdown complete")
}
