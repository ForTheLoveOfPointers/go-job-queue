package jobs

import "encoding/json"

type Status string

const (
	StatusPending   Status = "pending"
	StatusRunning   Status = "running"
	StatusCompleted Status = "completed"
	StatusFailed    Status = "failed"
)

type Job struct {
	ID      string
	Type    string
	Status  Status
	Payload json.RawMessage
	Result  any
	Error   string
}
