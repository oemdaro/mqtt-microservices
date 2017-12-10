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

	_ "github.com/joho/godotenv/autoload"
	"github.com/oemdaro/mqtt-microservices-example/data-service/appconfig"
	"github.com/oemdaro/mqtt-microservices-example/data-service/db"
	"github.com/oemdaro/mqtt-microservices-example/data-service/kafka"
	"github.com/oemdaro/mqtt-microservices-example/data-service/workque"
)

var (
	// cassandraPeers the cassandra host string
	cassandraPeers = flag.String("cassandras", os.Getenv("CASSANDRA_PEERS"), "The Cassandra database to connect to, as a comma separated list")
	// cassandraKeyspace the cassandra keyspace
	cassandraKeyspace = flag.String("keyspace", os.Getenv("CASSANDRA_KEYSPACE"), "The Cassandra keyspace")
	// kafkaPeers the kafka broker host string
	kafkaPeers = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	// kafkaGroupID the kafka consumer group id
	kafkaGroupID = flag.String("group-id", "data-service-group", "The Kafka consumer group id")
	// kafkaPeers the kafka topics string
	kafkaTopics = flag.String("topics", os.Getenv("KAFKA_TOPICS"), "The Kafka topics to subscribe to, as a comma separated list")
	// numMaxQueue max number of numMaxQueue
	numMaxQueue = flag.String("max-queue", os.Getenv("MAX_QUEUE"), "The maximum queues")
	// numMaxWorker max number of workers
	numMaxWorker = flag.String("max-worker", os.Getenv("MAX_WORKER"), "The maximum workers")
	// signals we want to gracefully shutdown when it receives a SIGTERM or SIGINT
	signals = make(chan os.Signal, 1)
	done    = make(chan bool, 1)
)

func main() {
	flag.Parse()

	if *cassandraPeers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *kafkaPeers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *numMaxWorker == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *numMaxQueue == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	maxWorker, err := strconv.Atoi(*numMaxWorker)
	if err != nil {
		log.Fatalf("Invalid MAX_WORKERS value: %s", err)
	}
	maxQueue, err := strconv.Atoi(*numMaxQueue)
	if err != nil {
		log.Fatalf("Invalid MAX_WORKERS value: %s", err)
	}
	log.Printf("Max Queue: %s, Max Worker: %s", *numMaxQueue, *numMaxWorker)

	cassandraList := strings.Split(*cassandraPeers, ",")
	log.Printf("Cassandra Peers: %s", strings.Join(cassandraList, ", "))
	log.Printf("Cassandra Keyspace: %s", *cassandraKeyspace)

	brokerList := strings.Split(*kafkaPeers, ",")
	log.Printf("Kafka Brokers: %s", strings.Join(brokerList, ", "))
	topics := strings.Split(*kafkaTopics, ",")
	log.Printf("Kafka Topics: %s", strings.Join(topics, ", "))

	// Load configuration
	appconfig.Load(cassandraList, *cassandraKeyspace, brokerList, *kafkaGroupID, topics)

	cassandraDB, err := db.NewDB()
	if err != nil {
		log.Fatalf("error occurred while try to connect to Cassandra %s", err)
	}
	defer cassandraDB.Close()

	consumer, err := kafka.NewConsumer()
	if err != nil {
		log.Fatalf("error occurred while try to connect to Kafka brokers %s", err)
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
					work := workque.Job{
						DB:      cassandraDB,
						Message: msg,
					}
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
