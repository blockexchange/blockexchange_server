const jwt = require('jsonwebtoken');
const logger = require("../logger");

module.exports = function(req, res, next) {
	var token = req.headers.authorization;
	try {
		const payload = jwt.verify(token, process.env.BLOCKEXCHANGE_KEY);
		req.claims = payload;
		// token valid
		next();
	} catch (e) {
		// not authenticated
		logger.error("token-middleware: unauthenticated", e);
		res.status(401).end();
	}
};
