package api

import (
	"encoding/json"
	"for-the-love-of-pointers/job-queue/internal/api/types"
	"for-the-love-of-pointers/job-queue/internal/api/utils"
	"net/http"
)

type JobService interface {
	CreateJob(req types.CreateJobRequest) (types.JobResponse, error)
	GetJob(id string) (types.JobResponse, error)
}

type Handler struct {
	jobs JobService
}

func NewHandler(jobs JobService) *Handler {
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

}
