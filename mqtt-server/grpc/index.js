const grpc = require('grpc')
const path = require('path')
const config = require('../config')

const PROTO_PATH = path.join(__dirname, '/../../pb/authclient.proto')
const proto = grpc.load(PROTO_PATH).pb

const client = new proto.Auth(config.grpc.host, grpc.credentials.createInsecure())
module.exports = client
