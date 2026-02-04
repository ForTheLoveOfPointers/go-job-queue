package main

import (
	"for-the-love-of-pointers/job-queue/internal/api"
	"for-the-love-of-pointers/job-queue/internal/jobs"
	"net/http"
)

var store *jobs.Store = jobs.NewStore()
var queue *jobs.Queue = jobs.NewQueue(25)

func main() {

	jobService := jobs.NewService(store, queue)
	handler := api.NewHandler(jobService)
	r := api.NewRouter(handler)
	http.ListenAndServe(":3333", r)

}
