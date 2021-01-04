const app = require("../app");
const schema_dao = require("../dao/schema");
const logger = require("../logger");

const { MANAGEMENT } = require("../permissions");
const tokenmiddleware = require("../middleware/token");
const permissioncheck = require("../middleware/permissioncheck");

app.delete("/api/schema/:id", tokenmiddleware, permissioncheck(MANAGEMENT), async function(req, res){
	logger.debug("DELETE /api/schema/:id", req.params.id);

	const schema = await schema_dao.get_by_id(req.params.id);
	if (schema.user_id != req.claims.user_id){
		res.status(403).end();
		return;
	}

	// delete schema
	await schema_dao.delete_by_id(req.params.id);
	res.end();
});
