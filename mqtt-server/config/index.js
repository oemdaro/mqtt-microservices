const env = require('env-var')

const sentinels = [
  {
    port: env.get('REDIS_SENTINELS_PORT_1', '26380').asIntPositive(),
    host: env.get('REDIS_SENTINELS_HOST_1', '127.0.0.1').asString()
  }, {
    port: env.get('REDIS_SENTINELS_PORT_2', '26381').asIntPositive(),
    host: env.get('REDIS_SENTINELS_HOST_2', '127.0.0.1').asString()
  }, {
    port: env.get('REDIS_SENTINELS_PORT_3', '26382').asIntPositive(),
    host: env.get('REDIS_SENTINELS_HOST_3', '127.0.0.1').asString()
  }
]

const cache = {
  sentinels: sentinels,
  name: 'mymaster',
  // password: env.get('REDIS_PASSWORD', 'mysecret').asString(),
  db: env.get('REDIS_CACHE_DB', '1').asIntPositive()
}

const persistence = {
  sentinels: sentinels,
  name: 'mymaster',
  // password: env.get('REDIS_PASSWORD', 'mysecret').asString(),
  db: env.get('REDIS_PERSISTENCE_DB', '3').asIntPositive(),
  maxSessionDelivery: 100, // maximum offline messages deliverable on client CONNECT, default is 1000
  packetTTL: function (packet) { // offline message TTL, default is disabled
    return 3600 // seconds
  }
}

const mq = {
  sentinels: sentinels,
  name: 'mymaster',
  // password: env.get('REDIS_PASSWORD', 'mysecret').asString(),
  db: env.get('REDIS_MQ_DB', '4').asIntPositive()
}

module.exports = {
  version: '1.0.0',
  port: env.get('PORT', '1883').asIntPositive(),
  concurrency: env.get('CONCURRENCY', '100').asIntPositive(),
  heartbeatInterval: env.get('HEARTBEAT', '60000').asIntPositive(), // milliseconds
  connectTimeout: env.get('TIMEOUT', '30000').asIntPositive(), // milliseconds
  grpc: {
    host: env.get('GRPC_SERVER', 'localhost:50051').asString()
  },
  redis: {
    cache: cache,
    persistence: persistence,
    mq: mq
  },
  kafka: {
    peers: env.get('KAFKA_PEERS', 'localhost:9092').asString(),
    topic: env.get('KAFKA_TOPIC', 'mqtt.data').asString(),
    connectTimeout: env.get('KAFKA_CONNECT_TIMEOUT', '10000').asIntPositive(),
    requestTimeout: env.get('KAFKA_REQUEST_TIMEOUT', '30000').asIntPositive(),
    autoConnect: env.get('KAFKA_AUTO_CONNECT', 'true').asBool()
  }
}
