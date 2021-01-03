const jwt = require('jsonwebtoken');

module.exports = function(user, permissions, options) {

	const payload = {
		username: user.name,
		user_id: user.id,
		type: user.type,
		mail: user.mail,
		permissions: permissions
	};

	return jwt.sign(payload, process.env.BLOCKEXCHANGE_KEY, options);
};
