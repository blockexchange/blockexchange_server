const jwt = require('jsonwebtoken');
const logger = require("../logger");

module.exports = function(authFn) {
  return function(req, res, next) {
    var token = req.headers.authorization;
    try {
      const payload = jwt.verify(token, process.env.BLOCKEXCHANGE_KEY);
      req.claims = payload;

      if (typeof(authFn) == "function"){
        if (authFn(payload)) {
          // ok
          next();
        } else {
          // unauthorized
          res.status(403).end();
        }
      } else {
        // no authorization check
        next();
      }
    } catch (e) {
      // not authenticated
			logger.error("token-middleware: unauthenticated", e);
      res.status(401).end();
    }
  };
};
