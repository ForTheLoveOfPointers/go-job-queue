package api

import (
	"encoding/json"
	"for-the-love-of-pointers/job-queue/internal/api/types"
	"for-the-love-of-pointers/job-queue/internal/api/utils"
	"for-the-love-of-pointers/job-queue/internal/jobs"
	"net/http"
	"strconv"
)

type JobService interface {
	CreateJob(req types.CreateJobRequest) (types.JobResponse, error)
	GetJob(id string) (types.JobResponse, error)
}

type Handler struct {
	jobs *jobs.Service
}

func NewHandler(jobs *jobs.Service) *Handler {
	return &Handler{jobs: jobs}
}

func (h *Handler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var req types.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Type == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	job, err := h.jobs.CreateJob(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, job)
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	job_id := r.URL.Query().Get("id")

	_, err := strconv.Atoi(job_id)
	if err != nil {
		http.Error(w, "id must be an integer", 400)
		return
	}

	job_res, err := h.jobs.GetJob(job_id)

	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	utils.WriteJSON(w, http.StatusFound, job_res)
}
