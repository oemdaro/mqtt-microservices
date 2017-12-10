package workque

import (
	"log"
)

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// JobQueue a buffered channel that we can send work requests on.
var JobQueue chan Job

// Dispatcher struct of workers channels
type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
	// maxWorkers maximum number of worker
	maxWorkers int
}

// NewDispatcher create new dispatcher
func NewDispatcher(maxWorkers int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkerPool: pool, maxWorkers: maxWorkers}
}

// Run start the workers
func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	log.Println("Worker queue dispatcher started...")
	for {
		select {
		case job := <-JobQueue:
			log.Printf("a dispatcher request received")
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
