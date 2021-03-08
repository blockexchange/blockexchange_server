const schema_dao = require("../dao/schema");

module.exports = function(param_id_field){
	return async function(req, res, next){
		const schema_id = req.params[param_id_field];
		if (!schema_id){
			res.status(403).json({ message: "schema_id not found" });
			return;
		}
		const schema = await schema_dao.get_by_id(schema_id);
		if (schema.user_id != req.claims.user_id){
			res.status(403).json({ message: "not the owner of the schema" });
			return;
		}

		// provide schema on the request
		req.schema = schema;

		next();
	};
};
