package workque

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

var quit = make(chan *sync.WaitGroup, 1)

// Payload a Kafka consumer message
type Payload struct {
	Message sarama.ConsumerMessage
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
		// register the current worker into the worker queue.
		w.WorkerPool <- w.JobChannel
		log.Println("Worker started...")

		for {
			select {
			case job := <-w.JobChannel:
				// we have received a work request.
				if err := job.Payload.InsertIntoDB(); err != nil {
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

// InsertIntoDB insert receive data into database
func (p *Payload) InsertIntoDB() error {
	log.Printf("worker process message - Topic: %s, Partition: %d, Offset: %d,\tMessageKey: %s,\tMessageValue: %s", p.Message.Topic, p.Message.Partition, p.Message.Offset, p.Message.Key, p.Message.Value)
	time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
	log.Println("work done")
	return nil
}
