const jwt = require('jsonwebtoken');

module.exports = function(user) {

	const payload = {
		username: user.name,
		user_id: user.id,
		role: user.role,
		type: user.type
	};

	//ROLES: "UPLOAD_ONLY", "MEMBER", "ADMIN"

	return jwt.sign(payload, process.env.BLOCKEXCHANGE_KEY);
};
