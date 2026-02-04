package jobs

import (
	"context"
	"fmt"
)

type WorkerPool struct {
	Queue   *Queue
	Workers int
}

func (w *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < w.Workers; i++ {
		go w.work(ctx, i)
	}
}

func (w *WorkerPool) work(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-w.Queue.ch:
			w.process(job, id)
		}
	}
}

func (w *WorkerPool) process(job *Job, id int) {
	switch job.Type {
	case "printer":
		fmt.Printf("Print job id: %s\n", job.ID)
	default:
		fmt.Println("No suitable processing function available yet")
	}
}
