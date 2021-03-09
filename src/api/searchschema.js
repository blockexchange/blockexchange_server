const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();
const logger = require("../logger");

const app = require("../app");
const schema_dao = require("../dao/schema");
const schemamod_dao = require("../dao/schemamod");
const user_dao = require("../dao/user");
const user_schema_star_dao = require("../dao/userschemastar");

async function enrich(schema){
	const schemamods = await schemamod_dao.find_all(schema.id);
	const user = await user_dao.get_by_id(schema.user_id);
	const schema_stars = await user_schema_star_dao.count_by_schema_id(schema.id);

	const mods = schemamods.map(sm => sm.mod_name);

	return Object.assign({}, schema, {
		search_tokens: null,
		mods: mods,
		user: {
			name: user.name
		},
		stars: schema_stars.stars
	});
}

// data='{"keywords": "description"}'
// curl -X POST 127.0.0.1:8080/api/searchschema --data "${data}" -H "Content-Type: application/json"
app.post('/api/searchschema', jsonParser, async function(req, res){
  logger.debug("POST /api/searchschema", req.body);

	let schemas = [];

	// select proper query
	if (req.body.user_id) {
		// just by user_id
		schemas = await schema_dao.find_by_user_id(req.body.user_id);

	} else if (req.body.user_name) {
		// by username
		schemas = await schema_dao.find_by_user_name(req.body.user_name);

	} else if (req.body.keywords) {
		// by keywords
		schemas = await schema_dao.find_by_keywords(req.body.keywords);

	} else if (req.body.schema_id) {
		// by unique id
		const schema = await schema_dao.get_by_id(req.body.schema_id);
		schemas.push(schema);

	} else if (req.body.schema_name && req.body.user_name) {
		// by schema name and user name (unique)
		const schema = await schema_dao.get_by_schemaname_and_username(req.body.schema_name, req.body.user_name);
		schemas.push(schema);
	}

	// TODO: sorting { sorting: { field: "created", desc: false }}
	// TODO: paging { paging: { page: 1, items: 20 }}

	const enriched_schemas = await Promise.all(schemas.map(enrich));
	res.json(enriched_schemas);
});

app.get("/api/searchrecent/:count", async function(req, res){
	const rows = await schema_dao.find_recent(req.params.count);
	const enriched_rows = await Promise.all(rows.map(enrich));
	res.json(enriched_rows);
});


app.get("/api/search/schema/byname/:username/:name", async function(req, res){
	logger.debug("POST /api/search/schema/byname/:username/:name", req.params);

	const schema = await schema_dao.get_by_schemaname_and_username(req.params.name, req.params.username);
	if (!schema) {
		res.status(204).end();
		return;
	}

	if (req.query.download === "true") {
    // increment download counter
    schema_dao.increment_downloads(schema.id);
  }

	const enriched_schema = await enrich(schema);
	res.json(enriched_schema);
});
