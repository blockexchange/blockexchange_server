const jwt = require('jsonwebtoken');

module.exports = function(user, options) {

	const payload = {
		username: user.name,
		user_id: user.id,
		role: user.role,
		type: user.type,
		mail: user.mail
	};

	//ROLES: "UPLOAD_ONLY", "MEMBER", "ADMIN"

	return jwt.sign(payload, process.env.BLOCKEXCHANGE_KEY, options);
};
