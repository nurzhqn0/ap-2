package queue

import "assignment2/internal/model"

// TaskQueue manages the queue of tasks to be processed
type TaskQueue struct {
	ch chan *model.Task
}

// NewTaskQueue creates a new task queue with specified buffer size
func NewTaskQueue(bufferSize int) *TaskQueue {
	return &TaskQueue{
		ch: make(chan *model.Task, bufferSize),
	}
}

// Enqueue adds a task to the queue
func (q *TaskQueue) Enqueue(task *model.Task) {
	q.ch <- task
}

// Dequeue returns the channel to receive tasks from
func (q *TaskQueue) Dequeue() <-chan *model.Task {
	return q.ch
}

// Close closes the queue channel
func (q *TaskQueue) Close() {
	close(q.ch)
}
