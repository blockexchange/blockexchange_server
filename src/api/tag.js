const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const tag_dao = require("../dao/tag");
const schematag_dao = require("../dao/schematag");
const schema_dao = require("../dao/schema");
const logger = require("../logger");

const { MANAGEMENT } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");


app.get('/api/tag', async function(req, res){
	logger.debug("GET /api/tag");

	const tags = await tag_dao.find_all();
	res.json(tags);
});

app.get('/api/schematag/:schema_id', function(req, res){
	logger.debug("GET /api/schematag/:schema_id");
	schematag_dao.find(req.params.schema_id).then(list => {
		res.json(list);
	});
});

app.put('/api/schematag/:schema_id/:tag_id', tokenmiddleware, permissioncheck(MANAGEMENT), jsonParser, async function(req, res){
	logger.debug("PUT /api/schematag/:schema_id/:tag_id");

	const schema = await schema_dao.get_by_id(req.params.schema_id);
	if (schema.user_id != req.claims.user_id){
		res.status(403).end();
		return;
	}

	await schematag_dao.create(req.params.schema_id, req.params.tag_id);
	res.end();
});

app.delete('/api/schematag/:schema_id/:tag_id', tokenmiddleware, permissioncheck(MANAGEMENT), jsonParser, async function(req, res){
	logger.debug("DELETE /api/schematag/:schema_id/:tag_id");

	const schema = await schema_dao.get_by_id(req.params.schema_id);
	if (schema.user_id != req.claims.user_id){
		res.status(403).end();
		return;
	}

	await schematag_dao.remove(req.params.schema_id, req.params.tag_id);
	res.end();
});
