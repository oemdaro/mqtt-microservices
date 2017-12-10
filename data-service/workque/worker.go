package workque

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gocql/gocql"
)

var quit = make(chan *sync.WaitGroup, 1)

// Data a MQTT payload model
type Data struct {
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

// Worker represents the worker that executes the job
type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
}

// NewWorker create new channel worker
func NewWorker(workerPool chan chan Job) Worker {
	return Worker{
		WorkerPool: workerPool,
		JobChannel: make(chan Job)}
}

// Start method starts the run loop for the worker,
// listening for a quit channel in case we need to stop it
func (w Worker) Start() {
	go func() {
		log.Println("Worker started")

		for {
			// register the current worker into the worker queue.
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.InsertIntoDB(); err != nil {
					log.Printf("Error inserting into database: %s", err.Error())
				}

			case wg := <-quit:
				// we have received a signal to stop
				log.Println("Worker stopped")
				wg.Done()
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w Worker) Stop(wg *sync.WaitGroup) {
	go func() {
		log.Println("Stopping worker...")
		quit <- wg
	}()
}

// InsertIntoDB insert received data into database
func (j *Job) InsertIntoDB() error {
	log.Printf("worker process message - Topic: %s, Partition: %d, Offset: %d,\tMessageKey: %s,\tMessageValue: %s", j.Message.Topic, j.Message.Partition, j.Message.Offset, j.Message.Key, j.Message.Value)
	var data Data
	if err := json.Unmarshal(j.Message.Value, &data); err != nil {
		return err
	}

	if err := j.DB.Query(`INSERT INTO data (id, username, temperature, humidity, timestamp) VALUES (?, ?, ?, ?, ?)`,
		gocql.TimeUUID(), j.Message.Key, data.Temperature, data.Humidity, time.Now()).Exec(); err != nil {
		return err
	}
	log.Println("work done")
	return nil
}
