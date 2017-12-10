package appconfig

type (
	// CassandraDB config
	CassandraDB struct {
		Peers    []string
		Keyspace string
	}
	// Kafka config
	Kafka struct {
		Peers   []string
		GroupID string
		Topics  []string
	}
	configuration struct {
		CassandraDB *CassandraDB
		Kafka       *Kafka
	}
)

// Config an application configuration
var Config *configuration

// Load the application configuration
func Load(cassandraPeers []string, keyspace string, kafkaPeers []string, kafkaGroupID string, kafkaTopics []string) {
	cassandra := &CassandraDB{
		Peers:    cassandraPeers,
		Keyspace: keyspace,
	}
	kafka := &Kafka{
		Peers:   kafkaPeers,
		GroupID: kafkaGroupID,
		Topics:  kafkaTopics,
	}

	Config = &configuration{
		CassandraDB: cassandra,
		Kafka:       kafka,
	}
}
