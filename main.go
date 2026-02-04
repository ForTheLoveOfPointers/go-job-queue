package main

import (
	"context"
	"for-the-love-of-pointers/job-queue/internal/api"
	"for-the-love-of-pointers/job-queue/internal/jobs"
	"net/http"
)

var store *jobs.Store = jobs.NewStore()
var queue *jobs.Queue = jobs.NewQueue(50)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wp := jobs.WorkerPool{Queue: queue, Workers: 5}

	wp.Start(ctx)

	/*
		START WORKER POOL HERE
	*/

	jobService := jobs.NewService(store, queue)
	handler := api.NewHandler(jobService)
	r := api.NewRouter(handler)
	http.ListenAndServe(":3333", r)

}
