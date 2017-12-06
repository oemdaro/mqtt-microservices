const grpc = require('../grpc')
const logger = require('../logger')
const redis = require('../redis')

/**
 *
 * @param {Object} client
 * @param {String} username
 * @param {String} password
 * @param {Function} callback
 */
module.exports = function (client, username, password, callback) {
  grpc.authClient({
    clientKey: client.id,
    username: username,
    clientSecret: password.toString()
  }, (err, response) => {
    if (err) {
      logger.warn({err: err}, err.message)
      callback(err, false)
      return
    }

    if (response.code !== '200') {
      let err = new Error(response.detail)
      let errStr = `failed to authenticate by using clientId '${client.id}' and username '${username}', server response error code '${response.code}'`
      logger.warn({err: err}, errStr)
      callback(err, false)
      return
    }

    redis.set('mqtt:client:' + client.id, username, (err) => {
      if (err) {
        logger.error({err: err}, err.message)
        callback(err, false)
        return
      }
      logger.info({clientId: client.id, username: username}, 'authentication success')
      callback(null, true)
    })
  })
}
