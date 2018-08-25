'use strict'
const merge = require('webpack-merge')
const prodEnv = require('./prod.env')

module.exports = merge(prodEnv, {
  NODE_ENV: '"development"',
  API_ENDPOINT: '"' + (process.env.API_ENDPOINT || 'http://localhost:5100') + '"'
})
