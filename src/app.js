const express = require('express');

const app = express();
app.use(express.static('public'));
app.disable('etag');

module.exports = app;
