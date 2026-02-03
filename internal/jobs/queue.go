package jobs

type Queue struct {
	ch chan *Job
}

func NewQueue(size int) *Queue {
	return &Queue{
		ch: make(chan *Job, size),
	}
}

func (q *Queue) Enqueue(job *Job) {
	q.ch <- job
}

func (q *Queue) Channel() <-chan *Job {
	return q.ch
}
