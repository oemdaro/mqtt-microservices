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

	"github.com/joho/godotenv"
	"gopkg.in/Shopify/sarama.v1"
)

var (
	// Brokers the kafka broker connection string
	Brokers *string
	// MaxQueue max number of queue
	MaxQueue *string
	// MaxWorker max number of workers
	MaxWorker *string
	// Verbose use to turn on Sarama logging
	Verbose *bool
	// signals we want to gracefully shutdown
	// when it receives a SIGTERM or SIGINT
	signals = make(chan os.Signal, 1)
	done    = make(chan bool, 1)
)

// Job represents the job to be run
type Job struct {
	Payload Payload
}

// Payload the coming data payload
type Payload struct {
}

// JobQueue a buffered channel that we can send work requests on.
var JobQueue chan Job

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Rename sample.env file to .env")
	}

	Brokers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	MaxQueue = flag.String("max-queue", os.Getenv("MAX_QUEUE"), "The maximum queues")
	MaxWorker = flag.String("max-worker", os.Getenv("MAX_WORKER"), "The maximum workers")
	Verbose = flag.Bool("verbose", false, "Turn on Sarama logging")
	flag.Parse()

	if *Verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}
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

	var wg sync.WaitGroup
	wg.Add(maxWorker)

	JobQueue = make(chan Job, maxQueue)
	dispatcher := NewDispatcher(maxWorker)
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
			case <-signals:
				log.Println("Graceful shutting down...")
				w := Worker{}
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
	log.Println("Successfully shutdown")
	os.Exit(0)
}
