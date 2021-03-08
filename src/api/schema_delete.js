const app = require("../app");
const schema_dao = require("../dao/schema");
const logger = require("../logger");

const { MANAGEMENT } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");
const schemaownercheck = require("../middleware/schemaownercheck");

app.delete("/api/schema/:id",
	tokenmiddleware,
	permissioncheck(MANAGEMENT),
	schemaownercheck("id"),
	async function(req, res){

	logger.debug("DELETE /api/schema/:id", req.params.id);

	// delete schema
	await schema_dao.remove(req.params.id);
	res.end();
});
