package worker

import (
	"context"
	"encoding/json"

	"github.com/KKKHEAO/task-processing/packages/domain"
	"github.com/KKKHEAO/task-processing/packages/kafka"
	"go.uber.org/zap"
)

const LIMIT = 100

type Job struct {
	Payload []byte
}

type Pool struct {
	ctx   context.Context
	jobs  chan Job
	retry *kafka.Producer
	dlq   *kafka.Producer
	log   *zap.Logger
}

func NewPool(ctx context.Context, maxWorkers int, retry *kafka.Producer, dlq *kafka.Producer, log *zap.Logger) *Pool {
	p := &Pool{
		ctx:   ctx,
		jobs:  make(chan Job),
		retry: retry,
		dlq:   dlq,
		log:   log,
	}

	for i := 0; i < maxWorkers; i++ {
		go p.Worker(i)
	}

	return p
}

func (p *Pool) Worker(id int) {
	for job := range p.jobs {
		p.log.Info("processing job", zap.Int("worker_id", id))

		if err := Process(job); err == nil {
			continue
		}

		var event domain.TaskCreatedEvent
		if err := json.Unmarshal(job.Payload, &event); err != nil {
			p.log.Error("failed to unmarshal event", zap.Error(err))
			continue
		}

		event.Retry++

		payload, err := json.Marshal(event)
		if err != nil {
			p.log.Error("failed to marshal event", zap.Error(err))
			continue
		}

		if event.Retry < 3 {
			p.log.Info("retry task", zap.String("task_id", event.Id))
			p.retry.Send(p.ctx, payload)
		} else {
			p.log.Info("send to dlq", zap.String("task_id", event.Id))
			p.dlq.Send(p.ctx, payload)
		}
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}
