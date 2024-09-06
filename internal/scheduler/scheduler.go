package scheduler

import (
	"errors"
	"log"
	"sync"

	"github.com/madchin/trader-bot/internal/gateway"
)

var ErrNoJobToSchedule = errors.New("no job to schedule")

type scheduler struct {
	mu *sync.Mutex
	q  *queue[gateway.Job]
}

// we ensure we only have one scheduler
var Scheduler = &scheduler{
	&sync.Mutex{},
	&queue[gateway.Job]{},
}

func (t *scheduler) Schedule(jobInfo gateway.Job) {
	t.mu.Lock()
	log.Printf("scheduling job %v", jobInfo)
	log.Printf("yo zimakl  %v", jobInfo)
	t.q.enqueue(jobInfo)
	t.mu.Unlock()
}

func (t *scheduler) Delegate() (gateway.Job, error) {
	t.mu.Lock()
	data, err := t.q.dequeue()
	t.mu.Unlock()
	if err != nil {
		return nil, ErrNoJobToSchedule
	}
	return data, nil
}

func (t *scheduler) Count() int {
	t.mu.Lock()
	size := t.q.size
	t.mu.Unlock()
	return size
}
