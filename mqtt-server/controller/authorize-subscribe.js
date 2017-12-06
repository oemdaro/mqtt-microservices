const logger = require('../logger')
const redis = require('../redis')

/**
 *
 * @param {Object} client MQTT client
 * @param {Object} subscription pattern
 * @param {Function} done callback function done(err,pattern)
 */
module.exports = (client, subscription, done) => {
  let topic = subscription.topic.split('/')
  if (topic.length < 2) {
    let errStr = `subscribe is not authorized, subscribe wrong topic '${subscription.topic}'`
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
      let errStr = `subscribe is not authorized, topic '${subscription.topic}' not match with client id '${client.id}'`
      logger.warn(errStr)
      done(new Error(errStr))
      return
    }

    logger.info(subscription, 'subscribe is authorized')
    done(null, subscription)
  })
}
