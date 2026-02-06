package types

import "encoding/json"

type CreateJobRequest struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type JobResponse struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
}
