const client = require('prom-client');

const collectDefaultMetrics = client.collectDefaultMetrics;
const registry = new client.Registry();

collectDefaultMetrics({ register: registry });

module.exports = registry;
