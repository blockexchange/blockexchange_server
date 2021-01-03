
const app = require("../app");
const access_token_dao = require("../dao/access_token");
const logger = require("../logger");

const bodyParser = require('body-parser');
const jsonParser = bodyParser.json({limit: '20mb'});

const { MANAGEMENT } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");

app.get('/api/access_token', tokenmiddleware, permissioncheck(MANAGEMENT), async function(req, res){
	logger.debug("GET /api/access_token");

	const user_id = req.claims.user_id;
	const tokens = await access_token_dao.find_all_by_userid(user_id);
	res.json(tokens);
});

app.delete('/api/access_token/:id', tokenmiddleware, permissioncheck(MANAGEMENT), async function(req, res){
	logger.debug("DELETE /api/access_token/:id", req.params.id);

	const user_id = req.claims.user_id;
	await access_token_dao.remove(user_id, +req.params.id);
	res.end();
});

app.post('/api/access_token', tokenmiddleware, permissioncheck(MANAGEMENT), jsonParser, async function(req, res){
	logger.debug("POST /api/access_token", req.body);

	if (!req.body.name){
		res.status(500).end();
		return;
	}

	const user_id = req.claims.user_id;
	const created = Date.now();
	const expires = req.body.expires;
	const name = req.body.name;
	const token = access_token_dao.generateToken();

	const access_token = await access_token_dao.create(user_id, created, expires, name, token);

	res.json(access_token);
});
