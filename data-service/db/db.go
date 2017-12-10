package db

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/oemdaro/mqtt-microservices-example/data-service/appconfig"
)

// DB is an instant hold database session
type DB struct {
	*gocql.Session
}

// NewDB return a new Cassandra session
func NewDB() (*DB, error) {
	config := appconfig.Config.CassandraDB
	cluster := gocql.NewCluster(config.Peers...)
	cluster.Keyspace = config.Keyspace
	cluster.Consistency = gocql.Quorum

	log.Println(config.Peers)
	log.Println(config.Keyspace)
	log.Println(gocql.Quorum)

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return &DB{session}, nil
}
