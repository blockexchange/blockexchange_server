var jwt = require('jsonwebtoken');

module.exports = function(req) {
  return new Promise((resolve, reject) => {
    var token = req.headers.authorization;
    try {
      const payload = jwt.verify(token, process.env.BLOCKEXCHANGE_KEY);
      resolve(payload);
    } catch (e) {
      reject(e);
    }
  });
};
