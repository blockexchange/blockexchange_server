var jwt = require('jsonwebtoken');

module.exports.rolecheck = function(rolename) {
  return function(req, res, next){
    if (!req.claims || req.claims.role != rolename){
      // unauthorized
      res.status(403).end();
      return;
    }
    next();
  };
};

module.exports.verifytoken = function() {
  return function(req, res, next) {
    var token = req.headers.authorization;
    try {
      const payload = jwt.verify(token, process.env.BLOCKEXCHANGE_KEY);
      req.claims = payload;

      // no authorization check
      next();
    } catch (e) {
      // not authenticated
      res.status(401).end();
    }
  };
};
