package worker

import (
	"context"
	"encoding/json"
	"log"
	"task-processing/internal/domain"
	"task-processing/internal/kafka"
)

const LIMIT = 100

type Job struct {
	Payload []byte
}

type Pool struct {
	jobs  chan Job
	retry *kafka.Producer
	dlq   *kafka.Producer
}

func NewPool(maxWorkers int, retry *kafka.Producer, dlq *kafka.Producer) *Pool {
	p := &Pool{
		jobs:  make(chan Job),
		retry: retry,
		dlq:   dlq,
	}

	for i := 0; i < maxWorkers; i++ {
		go p.Worker(i)
	}

	return p
}

func (p *Pool) Worker(id int) {
	for job := range p.jobs {
		log.Println("worker", id, "processing job")

		if err := Process(job); err == nil {
			continue
		}

		var event domain.TaskCreatedEvent
		json.Unmarshal(job.Payload, &event)

		event.Retry++

		payload, _ := json.Marshal(event)
		if event.Retry < 3 {
			log.Println("retry task", event.TaskId)
			p.retry.Send(context.Background(), payload)

		} else {
			log.Println("send to dlq", event.TaskId)
			p.dlq.Send(context.Background(), payload)
		}
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}
