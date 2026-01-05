package model

// Task represents a background task
type Task struct {
	ID      string `json:"id"`
	Payload string `json:"payload"`
	Status  string `json:"status"` // PENDING, IN_PROGRESS, DONE
}

// TaskRequest represents the json body for creating a task
type TaskRequest struct {
	Payload string `json:"payload"`
}

// Stats represents server statistics
type Stats struct {
	Submitted  int `json:"submitted"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
}
