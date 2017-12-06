const Redis = require('ioredis')
const config = require('../config')

const redis = new Redis(config.redis.cache)
module.exports = redis
