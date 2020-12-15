
const app = require("../app");
const schema_dao = require("../dao/schema");
const logger = require("../logger");

// curl 127.0.0.1:8080/api/schema/1
app.get('/api/schema/:id', async function(req, res){
	logger.debug("GET /api/schema/:id", req.params.id);

	if (req.query.download === "true") {
		// increment download counter
		schema_dao.increment_downloads(req.params.id);
	}

	try {
		const schema = await schema_dao.get_by_id(req.params.id);
		res.json(schema);
	} catch (e) {
		res.status(500).send(e.message);
	}
});
