package jobs

import (
	"context"
	"fmt"
	"sync"
)

type WorkerPool struct {
	Queue   *Queue
	Workers int
	Wg      sync.WaitGroup
}

func (w *WorkerPool) Start(ctx context.Context) {
	for i := 0; i < w.Workers; i++ {
		w.Wg.Add(1)
		go w.work(ctx, i)
	}
}

func (w *WorkerPool) work(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			w.Wg.Done()
			return
		case job, ok := <-w.Queue.ch:
			if !ok {
				w.Wg.Done()
				return
			}
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
