const app = require("../app");
const schema_dao = require("../dao/schema");
const logger = require("../logger");

const tokenmiddleware = require("../middleware/token");
const rolecheck = require("../util/rolecheck");
const tokencheck = tokenmiddleware(claims => rolecheck.can_delete(claims.role));


app.delete("/api/schema/:id", tokencheck, async function(req, res){
	logger.debug("DELETE /api/schema/:id", req.params.id, req.body);

	const schema = schema_dao.get_by_id(req.params.id);
	if (schema.user_id != req.claims.user_id){
		res.status(403).end();
		return;
	}

	// delete schema
	await schema_dao.delete_by_id(req.params.id);
	res.end();
});
