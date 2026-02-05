package main

import (
	"context"
	"fmt"
	"for-the-love-of-pointers/job-queue/internal/api"
	"for-the-love-of-pointers/job-queue/internal/jobs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var store *jobs.Store = jobs.NewStore()
var queue *jobs.Queue = jobs.NewQueue(50)

func main() {
	ctx, close := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer close()

	wp := jobs.WorkerPool{Queue: queue, Workers: 5, Wg: sync.WaitGroup{}}

	wp.Start(ctx)

	jobService := jobs.NewService(store, queue)
	handler := api.NewHandler(jobService)
	r := api.NewRouter(handler)

	server := &http.Server{
		Addr:    ":3333",
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()

	fmt.Println("\nGraceful shutdown started...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutdownCtx)

	wp.Wg.Wait()
	fmt.Println("DONE")
}
