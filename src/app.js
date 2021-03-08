const express = require('express');
const errorhandler = require("./middleware/errorhandler");

const app = express();
app.use("/node_modules", express.static('node_modules'));
app.use(express.static('public'));
app.use(errorhandler);
app.disable('etag');

module.exports = app;
