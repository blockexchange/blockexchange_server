const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const schemamod_dao = require("../dao/schemamod");
const schemagroup_dao = require("../dao/schemagroup");
const user_schema_star_dao = require("../dao/userschemastar");

// TODO: maybe optimize with subqueries
function enrich(schema){
	return Promise.all([
		schemamod_dao.find_all(schema.id),
		schemagroup_dao.get_by_id(schema.schemagroup_id),
		user_schema_star_dao.count_by_schema_id(schema.id)
	])
	.then(results => {
		const schemamods = results[0];
		const schemagroup = results[1];
		const stars = results[2];

		const mods = {};
		schemamods.forEach(mod => mods[mod.mod_name] = mod.node_count);

		return Object.assign({}, schema, {
			search_tokens: null,
			mods: mods,
			schemagroup: schemagroup,
			stars: +stars.stars
		});
	});
}

// data='{"keywords": "description"}'
// curl -X POST 127.0.0.1:8080/api/searchschema --data "${data}" -H "Content-Type: application/json"
app.post('/api/searchschema', jsonParser, function(req, res){
  console.log("POST /api/searchschema", req.body);

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

	} else if (req.body.schema_name && req.body.group_name) {
		// by schema name and group_name
		q = schema_dao.get_by_schemaname_and_groupname(req.body.schema_name, req.body.group_name)
		.then(schema => [schema]);

	} else if (req.body.group_name) {
		// by group_name
		q = schema_dao.find_by_schemagroup_name(req.body.group_name);

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

app.get("/api/searchrecent/:count", function(req, res){
	schema_dao.find_recent(req.params.count)
	.then(rows => {
		return Promise.all(rows.map(enrich))
		.then(enriched_rows => res.json(enriched_rows));
	})
	.catch(e => {
		console.error(e);
		res.status(500).end();
	});

});


app.get("/api/search/schema/byname/:group_name/:name", function(req, res){
	console.log("POST /api/search/schema/byname/:group_name/:name", req.params);

	schema_dao.get_by_schemaname_and_groupname(req.params.name, req.params.group_name)
	.then(schema => {
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
