const { createLogger, format, transports } = require('winston');


module.exports = createLogger({
  level: process.env.LOGLEVEL || "info",
  format: format.combine(
    format.timestamp({
      format: 'YYYY-MM-DD HH:mm:ss'
    }),
    format.printf(info => `${info.timestamp} ${info.level}: ${info.message}`)
  ),
  transports: [
    new transports.Console()
  ]
});
