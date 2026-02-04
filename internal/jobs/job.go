package jobs

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
	Payload []byte
	Result  any
	Error   string
}
