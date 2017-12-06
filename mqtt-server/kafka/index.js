const kafka = require('kafka-node')
const config = require('../config')

const HighLevelProducer = kafka.HighLevelProducer
const client = new kafka.KafkaClient({
  kafkaHost: config.kafka.peers,
  connectTimeout: config.kafka.connectTimeout,
  requestTimeout: config.kafka.requestTimeout,
  autoConnect: config.kafka.autoConnect
})
const producer = new HighLevelProducer(client)

module.exports = producer
