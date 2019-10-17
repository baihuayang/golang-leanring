package internal

import "sync"

type Job func() error

type AsyncJobQueue struct {
	lock sync.Locker
	jobs []Job
}

func NewAsyncJobQueue() AsyncJobQueue {
	return AsyncJobQueue{
		lock: SpinLock(),
	}
}

func (q *AsyncJobQueue) Push(job Job) {
	q.lock.Lock()
	q.jobs = append(q.jobs, job)
	q.lock.Unlock()
}

func (q *AsyncJobQueue) Walk() error {
	q.lock.Lock()
	jobs := q.jobs
	q.jobs = nil
	q.lock.Unlock()
	for _, job := range jobs {
		if err := job(); err != nil {
			return err
		}
	}
	return nil
}
