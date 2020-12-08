const jwt = require('jsonwebtoken');

module.exports = function(user) {

	const payload = {
		username: user.name,
		user_id: user.id
	};

	// check for temporary role, only allow creation of content with it
	if (user.role == "TEMP"){
		// temporary/default user
		payload.permissions = {
			schema: {
				create: true
			}
		};
	} else {
		// normal user
		payload.permissions = {
			user: {
				update: true
			},
			schema: {
				create: true,
				update: true,
				delete: true
			},
			screenshot: {
				create: true,
				delete: true
			}
		};
	}
	// TODO: custom permissions on token for normal users
	return jwt.sign(payload, process.env.BLOCKEXCHANGE_KEY);
};
