
const app = require("../app");
const collection_schema_dao = require("../dao/collection_schema");
const collection_dao = require("../dao/collection");
const logger = require("../logger");

const bodyParser = require('body-parser');
const jsonParser = bodyParser.json({limit: '20mb'});

const { MANAGEMENT, UPLOAD } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");

app.get('/api/collection_schema/by-collectionid/:collection_id', async function(req, res){
	logger.debug("GET /api/collection_schema/by-collectionid", req.params.collection_id);

	const collection_schemas = await collection_schema_dao.find_all_by_collectionid(+req.params.collection_id);
	res.json(collection_schemas);
});

app.post('/api/collection_schema', tokenmiddleware, permissioncheck(UPLOAD), jsonParser, async function(req, res){
	logger.debug("POST /api/collection_schema", req.body);

	const user_id = req.claims.user_id;

	const collection = await collection_dao.find_by_collectionid(req.body.collection_id);
	if (!collection){
		return res.status(403).end();
	}

	if (collection.user_id != user_id){
		return res.status(403).end();
	}

	const collection_schema = await collection_schema_dao.create(req.body);
	res.json(collection_schema);
});

app.delete('/api/collection_schema/:collection_id/:schema_id', tokenmiddleware, permissioncheck(MANAGEMENT), jsonParser, async function(req, res){
	logger.debug("DELETE /api/collection_schema/:collection_id/:schema_id", req.params.collection_id, req.params.schema_id);

	const user_id = req.claims.user_id;

	const collection = await collection_dao.find_by_collectionid(req.params.collection_id);
	if (!collection){
		return res.status(403).end();
	}

	if (collection.user_id != user_id){
		return res.status(403).end();
	}

	await collection_schema_dao.remove(req.params.collection_id, req.params.schema_id);
	res.end();
});
