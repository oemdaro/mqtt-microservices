const { createLogger, stdSerializers } = require('bunyan')

const logger = createLogger({
  name: 'mqtt-server',
  serializers: {
    err: stdSerializers.err
  }
})

module.exports = logger
