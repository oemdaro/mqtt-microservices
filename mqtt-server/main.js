require('dotenv').config()
const aedes = require('aedes')
const redisMQ = require('mqemitter-redis')
const persistence = require('aedes-persistence-redis')
const { series } = require('async')
const config = require('./config')
const producer = require('./kafka')
const redis = require('./redis')
const controller = require('./controller')
const logger = require('./logger')

const aedesMQ = redisMQ(config.redis.mq)
const aedesPersistence = persistence(config.redis.persistence)
const mqtt = aedes({
  mq: aedesMQ,
  persistence: aedesPersistence,
  concurrency: config.concurrency,
  heartbeatInterval: config.heartbeatInterval,
  connectTimeout: config.connectTimeout,
  authenticate: controller.authenticate,
  authorizePublish: controller.authorizePublish,
  authorizeSubscribe: controller.authorizeSubscribe,
  published: controller.published
})
const server = require('net').createServer(mqtt.handle)

// emits when Kafka Producer is ready
producer.on('ready', () => {
  logger.info('kafka Producer is ready')
  server.listen(config.port, () => {
    logger.info(`server listening on port ${config.port}`)
  })
})
// emits when an error occurs in Kafka Producer
producer.on('error', (err) => {
  logger.error({err: err}, err.message)
})

// emits when a connection is established to the Redis server
redis.on('connect', () => {
  logger.error('connection is established to the Redis server')
})
// emits when an error occurs while connecting
redis.on('error', (err) => {
  logger.error(err.message)
})

// emits when a Client disconnects
mqtt.on('clientDisconnect', (client) => {
  logger.info(`remove Redis record for client '${client.id}'`)
  redis.del('mqtt:client:' + client.id)
})

let exit = () => {
  logger.info('Graceful shutting down...')
  series([
    function (callback) {
      server.close(() => callback())
    },
    function (callback) {
      mqtt.close(() => callback())
    },
    function (callback) {
      producer.close(() => callback())
    },
    function (callback) {
      redis.quit((err, res) => callback(err, res))
    }
  ], (err, result) => {
    if (err) {
      logger.error({err: err}, err.message)
    } else {
      logger.info('Shut down completed')
      process.exit(0)
    }
  })
}

process.on('SIGINT', () => exit())
process.on('SIGTERM', () => exit())
