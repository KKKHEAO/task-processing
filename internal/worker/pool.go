package worker

import "log"

const LIMIT = 100

type Job struct {
	Payload []byte
}

type Pool struct {
	jobs chan Job
}

func NewPool(maxWorkers int) *Pool {
	p := &Pool{
		jobs: make(chan Job),
	}

	for i := 0; i < maxWorkers; i++ {
		go p.Worker(i)
	}

	return p
}

func (p *Pool) Worker(id int) {
	for job := range p.jobs {
		log.Println("worker", id, "processing job")

		Process(job)
	}
}

func (p *Pool) Submit(job Job) {
	p.jobs <- job
}
