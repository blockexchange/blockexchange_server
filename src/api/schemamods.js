const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();
const logger = require("../logger");

const app = require("../app");
const schemamod_dao = require("../dao/schemamod");

const { UPLOAD } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");
const schemaownercheck = require("../middleware/schemaownercheck");

app.get('/api/schema/:id/mods', async function(req, res){
	logger.debug("GET /api/schema/:id/mods", req.params.id);

	const schemamods = await schemamod_dao.find_all(req.params.id);
	res.json(schemamods.map(sm => sm.mod_name));
});


app.post('/api/schema/:id/mods',
	tokenmiddleware,
	permissioncheck(UPLOAD),
	schemaownercheck("id"),
	jsonParser,
	async function(req, res){

	logger.debug("POST /api/schema/:id/mods", req.params.id, req.body);

	req.body.forEach(mod_name => {
		schemamod_dao.create(req.params.id, mod_name);
	});

	res.end();
});
