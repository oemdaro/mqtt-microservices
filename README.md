# mqtt-microservices

The example of MQTT microservices based system build on top of Node.js, Go (Golang), gRPC and Apache Kafka

## Dockerize App

To build docker images of our app run the following command

```bash
$ cd $GOPATH/src/github.com/oemdaro/mqtt-microservices
$ make
```

## Run Docker App

```bash
$ # start data-service app
$ docker run -t -i --rm --network="host" --env CASSANDRA_PEERS=127.0.0.1,127.0.0.2,127.0.0.3 --env CASSANDRA_KEYSPACE=mqttexample --env KAFKA_PEERS=localhost:9092 --env KAFKA_TOPICS=mqtt.data --env MAX_QUEUE=5 --env MAX_WORKER=3 local/data-service
$ # start auth-service
$ docker run -t -i --rm --network="host" --env MYSQL_HOST=localhost:3306 --env MYSQL_DB=mqtt_example --env MYSQL_USER=mqttuser --env MYSQL_PASSWORD=password --env AES_KEY="a very very very very secret key" local/auth-service
$ # start mqtt-server
$ docker run -t -i --rm --network="host" --env PORT=1883 --env REDIS_SENTINELS_HOST_1=127.0.0.1 --env REDIS_SENTINELS_PORT_1=26380 --env REDIS_SENTINELS_HOST_2=127.0.0.1 --env REDIS_SENTINELS_PORT_2=26381 --env REDIS_SENTINELS_HOST_3=127.0.0.1 --env REDIS_SENTINELS_PORT_3=26382 --env KAFKA_PEERS=localhost:9092 --env KAFKA_TOPIC=mqtt.data --env GRPC_SERVER=localhost:50051 local/mqtt-server
```

## Run on Kubernetes

> Update soon...
