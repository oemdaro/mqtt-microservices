package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	cluster "github.com/bsm/sarama-cluster"
	_ "github.com/joho/godotenv/autoload"
	"github.com/oemdaro/mqtt-microservices-example/data-service/workque"
)

var (
	// Brokers the kafka broker connection string
	Brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	// MaxQueue max number of queue
	MaxQueue = flag.String("max-queue", os.Getenv("MAX_QUEUE"), "The maximum queues")
	// MaxWorker max number of workers
	MaxWorker = flag.String("max-worker", os.Getenv("MAX_WORKER"), "The maximum workers")
	// signals we want to gracefully shutdown when it receives a SIGTERM or SIGINT
	signals = make(chan os.Signal, 1)
	done    = make(chan bool, 1)
)

func main() {
	flag.Parse()

	if *Brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *MaxWorker == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *MaxQueue == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	maxWorker, err := strconv.Atoi(*MaxWorker)
	if err != nil {
		log.Printf("Invalid MAX_WORKERS value: %s", err)
	}
	maxQueue, err := strconv.Atoi(*MaxWorker)
	if err != nil {
		log.Printf("Invalid MAX_WORKERS value: %s", err)
	}

	brokerList := strings.Split(*Brokers, ",")
	log.Printf("Kafka Brokers: %s", strings.Join(brokerList, ", "))
	log.Printf("Max Worker: %s", *MaxWorker)

	// init (custom) config, enable errors and notifications
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true

	// init consumer
	topics := []string{"mqtt.data"}
	consumer, err := cluster.NewConsumer(brokerList, "data-service-group", topics, config)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("error occurred while closing consumer %s", err)
		}
	}()

	// consume errors
	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	// consume notifications
	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(maxWorker)

	workque.JobQueue = make(chan workque.Job, maxQueue)
	dispatcher := workque.NewDispatcher(maxWorker)
	dispatcher.Run()

	// Notify when receive SIGINT or SIGTERM
	// kill -SIGINT <PID> or Ctrl+c
	// kill -SIGTERM <PID>
	signal.Notify(signals,
		syscall.SIGINT,
		syscall.SIGTERM)

	go func() {
		for {
			select {
			case msg, ok := <-consumer.Messages():
				if ok {
					// let's create a job with the message
					work := workque.Job{Payload: workque.Payload{Message: *msg}}
					log.Println("sending message to workque")
					// Push the work onto the queue.
					workque.JobQueue <- work
					log.Println("sent message to workque")
					consumer.MarkOffset(msg, "") // mark message as processed
				}
			case <-signals:
				log.Println("Graceful shutting down...")
				w := workque.Worker{}
				for i := 0; i < maxWorker; i++ {
					w.Stop(&wg)
				}
				done <- true
			}
		}
	}()

	// Exiting
	<-done
	wg.Wait()
	log.Println("Shut down completed")
	os.Exit(0)
}
