
const app = require("../app");
const collection_dao = require("../dao/collection");
const logger = require("../logger");

const bodyParser = require('body-parser');
const jsonParser = bodyParser.json({limit: '20mb'});

const { MANAGEMENT } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");

app.get('/api/collection/by-userid/:user_id', async function(req, res){
	logger.debug("GET /api/collection/by-userid", req.params.user_id);

	const collections = await collection_dao.find_all_by_userid(+req.params.user_id);
	res.json(collections);
});

app.delete('/api/collection/:id', tokenmiddleware, permissioncheck(MANAGEMENT), async function(req, res){
	logger.debug("DELETE /api/collection/:id", req.params.id);

	const user_id = req.claims.user_id;
	await collection_dao.remove(user_id, +req.params.id);
	res.end();
});

app.post('/api/collection', tokenmiddleware, permissioncheck(MANAGEMENT), jsonParser, async function(req, res){
	logger.debug("POST /api/collection", req.body);

	if (!req.body.name){
		res.status(500).end();
		return;
	}

	const user_id = req.claims.user_id;
	const name = req.body.name;

	const collection = await collection_dao.create({
		user_id: user_id,
		name: name
	});

	res.json(collection);
});

app.post('/api/collection/:id', tokenmiddleware, permissioncheck(MANAGEMENT), jsonParser, async function(req, res){
	logger.debug("POST /api/collection/:id", req.body);

	if (!req.body.name || req.body.id){
		res.status(500).end();
		return;
	}

	const user_id = req.claims.user_id;
	const id = +req.body.id;
	const name = req.body.name;

	const collection = await collection_dao.update({
		user_id: user_id,
		id: id,
		name: name
	});

	res.json(collection);
});
