package worker

import (
	"encoding/json"
	"log"
	"task-processing/internal/domain"
	"time"
)

func Process(job Job) {

	var event domain.TaskCreatedEvent

	err := json.Unmarshal(job.Payload, &event)
	if err != nil {
		log.Println("json error:", err)
		return
	}

	log.Println("processing task:", event.TaskId)

	// имитация работы
	time.Sleep(20 * time.Second)

	log.Println("task done:", event.TaskId)
}
