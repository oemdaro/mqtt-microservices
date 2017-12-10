const producer = require('../kafka')
const logger = require('../logger')

/**
 * Send a valid packet into kafka for later process
 *
 * @param {Buffer} packet the published MQTT packet
 * @param {Object} client client will be null for internal messages
 * @param {Function} done callback function
 */
module.exports = (packet, client, done) => {
  if (client == null) {
    logger.debug(packet, 'internal message')
    done()
    return
  }

  let payloads = [
    {topic: 'mqtt.data', key: client.id, messages: packet.payload}
  ]
  producer.send(payloads, (err, data) => {
    if (err) {
      logger.error({err: err}, 'an error occurred when send data to kafka')
      done()
      return
    }

    logger.info('success send data to kafka')
    done()
  })
}
