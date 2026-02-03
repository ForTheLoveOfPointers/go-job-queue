package types

type CreateJobRequest struct {
	Type    string `json:"type"`
	Payload []byte `json:"payload"`
}

type JobResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}
