package workers

import (
	"context"
	"log"
	"sync"
	"time"

	"task_manager/Domain"
	"task_manager/Usecases"
)

type WorkerPool struct {
	n          int
	uc         Usecases.ReminderUsecase
	jobTimeout time.Duration
	wg         sync.WaitGroup
}

func NewWorkerPool(n int, uc Usecases.ReminderUsecase, jobTimeout time.Duration) *WorkerPool {
	return &WorkerPool{n: n, uc: uc, jobTimeout: jobTimeout}
}

// Start launches n workers that run until the jobs channel is closed and drained.
func (p *WorkerPool) Start(jobs <-chan Domain.Task) {
	for i := 1; i <= p.n; i++ {
		p.wg.Add(1)
		go p.worker(i, jobs)
	}
}

func (p *WorkerPool) worker(id int, jobs <-chan Domain.Task) {
	defer p.wg.Done()
	for task := range jobs {
		// Deliberately from context.Background(), NOT the shutdown context,
		// so a job already in progress can finish during shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), p.jobTimeout)
		log.Printf("worker %d: delivering task %v", id, task.ID)
		if err := p.uc.Deliver(ctx, task); err != nil {
			log.Printf("worker %d: task %v failed: %v", id, task.ID, err)
		}
		cancel()
	}
	log.Printf("worker %d: jobs channel closed, exiting", id)
}

// Wait blocks until every worker has exited.
func (p *WorkerPool) Wait() { p.wg.Wait() }
