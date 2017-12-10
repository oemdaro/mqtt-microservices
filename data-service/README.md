# data-service

The service that subscribe to Kafka topic and then insert received message into Cassandra database.

```bash
$ # use connect to Cassandra using cqlsh and then create Keyspace and Table
$ cqlsh
cqlsh> create keyspace mqttexample with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 2 };
cqlsh> create table mqttexample.data(id uuid PRIMARY KEY, username text, temperature text, humidity text, timestamp timestamp);
cqlsh> create index on mqttexample.data(username);
cqlsh> select * from mqttexample.data;

 id | humidity | temperature | timestamp | username
----+----------+-------------+-----------+----------

(0 rows)
cqlsh> exit
```
