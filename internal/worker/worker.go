package worker

import (
	"assignment2/internal/model"
	"assignment2/internal/queue"
	"assignment2/internal/store"
	"log"
	"sync"
	"time"
)

// Pool manages a pool of workers
type Pool struct {
	taskQueue  *queue.TaskQueue
	repository *store.Repository[string, *model.Task]
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewPool creates a new worker pool
func NewPool(taskQueue *queue.TaskQueue, repository *store.Repository[string, *model.Task]) *Pool {
	return &Pool{
		taskQueue:  taskQueue,
		repository: repository,
		stopCh:     make(chan struct{}),
	}
}

// Start starts the worker pool with specified number of workers
func (p *Pool) Start(numWorkers int) {
	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	// start monitoring goroutine (Part C)
	p.wg.Add(1)
	go p.monitor()
}

// worker processes tasks from the queue
func (p *Pool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case task, ok := <-p.taskQueue.Dequeue():
			if !ok {
				log.Printf("Worker %d: queue closed, stopping", id)
				return
			}

			log.Printf("Worker %d: processing task %s", id, task.ID)

			task.Status = "IN_PROGRESS"
			p.repository.Set(task.ID, task)

			// just simulation
			time.Sleep(2 * time.Second)

			// status to done
			task.Status = "DONE"
			p.repository.Set(task.ID, task)

			log.Printf("Worker %d: completed task %s", id, task.ID)

		case <-p.stopCh:
			log.Printf("Worker %d: received stop signal", id)
			return
		}
	}
}

// monitor logs statistics
func (p *Pool) monitor() {
	defer p.wg.Done()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			stats := p.getStats()
			log.Printf("MONITOR: Submitted=%d, InProgress=%d, Completed=%d",
				stats.Submitted, stats.InProgress, stats.Completed)

		case <-p.stopCh:
			log.Println("Monitor: received stop signal")
			return
		}
	}
}

// getStats calculates current statistics
func (p *Pool) getStats() model.Stats {
	tasks := p.repository.GetAll()

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

	return stats
}

// Stop stops all workers
func (p *Pool) Stop() {
	log.Println("Stopping worker pool...")
	close(p.stopCh)
	p.taskQueue.Close()
	p.wg.Wait()
	log.Println("All workers stopped")
}
