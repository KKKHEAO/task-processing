package worker

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"github.com/KKKHEAO/task-processing/packages/domain"
	"time"
)

func Process(job Job) error {

	var event domain.TaskCreatedEvent

	err := json.Unmarshal(job.Payload, &event)
	if err != nil {
		//log.Println("json error:", err)
		return err
	}

	log.Println("processing task:", event.TaskId)

	if rand.Intn(2) == 0 {
		return errors.New("random error")
	}

	// имитация работы
	time.Sleep(20 * time.Second)

	log.Println("task done:", event.TaskId)
	return nil
}
