const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();
const logger = require("../logger");

const app = require("../app");
const schema_dao = require("../dao/schema");
const schemamod_dao = require("../dao/schemamod");
const user_dao = require("../dao/user");
const user_schema_star_dao = require("../dao/userschemastar");

// TODO: maybe optimize with subqueries
function enrich(schema){
	return Promise.all([
		schemamod_dao.find_all(schema.id),
		user_dao.get_by_id(schema.user_id),
		user_schema_star_dao.count_by_schema_id(schema.id)
	])
	.then(results => {
		const mods = {};
		results[0].forEach(mod => mods[mod.mod_name] = mod.node_count);

		return Object.assign({}, schema, {
			search_tokens: null,
			mods: mods,
			user: {
				name: results[1].name
			},
			stars: +results[2].stars
		});
	});
}

// data='{"keywords": "description"}'
// curl -X POST 127.0.0.1:8080/api/searchschema --data "${data}" -H "Content-Type: application/json"
app.post('/api/searchschema', jsonParser, function(req, res){
  logger.debug("POST /api/searchschema", req.body);

	var q;

	// select proper query
	if (req.body.user_id) {
		// just by user_id
		q = schema_dao.find_by_user_id(req.body.user_id);

	} else if (req.body.user_name) {
		// by username
		q = schema_dao.find_by_user_name(req.body.user_name);

	} else if (req.body.keywords) {
		// by keywords
		q = schema_dao.find_by_keywords(req.body.keywords);

	} else if (req.body.schema_id) {
		// by unique id
		q = schema_dao.get_by_id(req.body.schema_id).then(schema => [schema]);

	} else if (req.body.schema_name && req.body.user_name) {
		// by schema name and user name (unique)
		q = schema_dao.get_by_schemaname_and_username(req.body.schema_name, req.body.user_name)
		.then(schema => [schema]);

	} else {
		// nothing to search for
		res.json([]);
		return;
	}

	// TODO: sorting { sorting: { field: "created", desc: false }}
	// TODO: paging { paging: { page: 1, items: 20 }}

  q.then(rows => {
		// enrich with additional data and return
		return Promise.all(rows.map(enrich))
		.then(enriched_rows => res.json(enriched_rows));
	})
  .catch(e => {
    console.error(e);
    res.status(500).end();
  });

});

app.get("/api/searchrecent/:count", async function(req, res){
	const rows = await schema_dao.find_recent(req.params.count);
	const enriched_rows = await Promise.all(rows.map(enrich));
	res.json(enriched_rows);
});


app.get("/api/search/schema/byname/:username/:name", function(req, res){
	logger.debug("POST /api/search/schema/byname/:username/:name", req.params);

	schema_dao.get_by_schemaname_and_username(req.params.name, req.params.username)
	.then(schema => {
		console.log(req.params, schema);
		if (!schema) {
			res.status(404).end();
			return;
		}

		if (req.query.download === "true") {
	    // increment download counter
	    schema_dao.increment_downloads(schema.id);
	  }

		return enrich(schema)
		.then(enriched_schema => {
			res.json(enriched_schema);
		});
	})
	.catch(e => {
		console.error(e);
		res.status(500).end();
	});

});
