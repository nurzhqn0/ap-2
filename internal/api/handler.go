package api

import (
	"assignment2/internal/model"
	"assignment2/internal/queue"
	"assignment2/internal/store"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

// Handler manages http requests
type Handler struct {
	repository    *store.Repository[string, *model.Task]
	taskQueue     *queue.TaskQueue
	taskIDCounter atomic.Uint64
}

// NewHandler creates a new api handler
func NewHandler(repository *store.Repository[string, *model.Task], taskQueue *queue.TaskQueue) *Handler {
	return &Handler{
		repository: repository,
		taskQueue:  taskQueue,
	}
}

// CreateTask handles post /tasks
func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// generates unique id
	taskID := fmt.Sprintf("%d", h.taskIDCounter.Add(1))

	// task with pending status
	task := &model.Task{
		ID:      taskID,
		Payload: req.Payload,
		Status:  "PENDING",
	}

	h.repository.Set(taskID, task)

	// enqueue for processing
	h.taskQueue.Enqueue(task)

	// 201 - created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTasks handles GET /tasks
func (h *Handler) GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks := h.repository.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// GetTask handles GET /tasks/{id}
func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// extract id from path
	taskID := r.URL.Path[len("/tasks/"):]

	task, exists := h.repository.Get(taskID)
	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// GetStats handles GET /stats
func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tasks := h.repository.GetAll()

	stats := model.Stats{}
	for _, task := range tasks {
		stats.Submitted++

		switch task.Status {
		case "IN_PROGRESS":
			stats.InProgress++
		case "DONE":
			stats.Completed++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
