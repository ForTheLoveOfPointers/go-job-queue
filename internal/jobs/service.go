package jobs

import (
	"errors"
	"for-the-love-of-pointers/job-queue/internal/api/types"

	"github.com/google/uuid"
)

type Service struct {
	store *Store
	queue *Queue
}

func NewService(store *Store, queue *Queue) *Service {
	return &Service{store: store, queue: queue}
}

func (s *Service) CreateJob(req types.CreateJobRequest) (*Job, error) {

	job := &Job{
		ID:      uuid.NewString(),
		Type:    req.Type,
		Status:  StatusPending,
		Payload: req.Payload,
	}

	s.store.Save(job)
	s.queue.Enqueue(job)

	return job, nil
}

func (s *Service) GetJob(id string) (*Job, error) {
	job, ok := s.store.Get(id)
	if !ok {
		return nil, errors.New("job not found")
	}
	return job, nil
}
