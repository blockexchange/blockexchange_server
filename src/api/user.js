const app = require("../app");
const user_dao = require("../dao/user");
const logger = require("../logger");


app.get('/api/user', async function(req, res){
	logger.debug("GET /api/user");

	const users = await user_dao.get_all();
	const list = users.map(user => {
		return {
			name: user.name,
			id: user.id,
			type: user.type,
			role: user.role,
			created: user.created
		};
	});

	res.json(list);
});
