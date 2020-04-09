const express = require('express');
const promMid = require('express-prometheus-middleware');

const app = express();
app.use(express.static('public'));
app.disable('etag');
app.use(promMid({
  metricsPath: '/metrics',
  collectDefaultMetrics: true,
  requestDurationBuckets: [0.01, 0.1, 0.25, 0.5, 1.0],
}));

module.exports = app;
