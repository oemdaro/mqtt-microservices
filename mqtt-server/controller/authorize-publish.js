const Joi = require('joi')
const logger = require('../logger')
const redis = require('../redis')

const schema = Joi.object().keys({
  temperature: Joi.string().regex(/^[0-9.]{1,15}$/).required(),
  humidity: Joi.string().regex(/^[0-9.]{1,15}$/).required()
})

/**
 *
 * @param {Object} client MQTT client
 * @param {Object} packet the packet publish by client
 * @param {Function} done callback function done(err)
 */
module.exports = (client, packet, done) => {
  let topic = packet.topic.split('/')
  if (topic.length < 2) {
    let errStr = `publish is not authorized, publish wrong topic '${packet.topic}'`
    logger.warn(errStr)
    done(new Error(errStr))
    return
  }

  redis.get('mqtt:client:' + client.id, (err, res) => {
    if (err) {
      logger.error({err: err}, err.message)
      done(err)
      return
    }

    if (topic[0] !== res) {
      let errStr = `publish is not authorized, topic '${packet.topic}' not match with client id '${client.id}'`
      logger.warn(errStr)
      done(new Error(errStr))
      return
    }

    Joi.validate(packet.payload.toString(), schema, (err, value) => {
      if (err) {
        logger.warn({payload: packet.payload.toString()}, `${client.id} publish invalid payload`)
        done(err)
        return
      }

      logger.info('publish is authorized')
      done(null)
    })
  })
}
