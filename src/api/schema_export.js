
const app = require("../app");
const schema_dao = require("../dao/schema");
const schemapart_dao = require("../dao/schemapart");
const logger = require("../logger");
const serializer = require("../util/serializer");
const export_part = require("../util/export_part");


app.get('/api/export/:id/:name', async function(req, res){
	logger.debug("GET /api/export/:id/:name", req.params.id, req.params.name);

	try {
		const schema = await schema_dao.get_by_id(req.params.id);
		if (schema.total_parts > 50){
			return res.status(500).send("WE export is limited to 50 parts!");
		}
		// open table
		let mts = "5:return {";

		for (let x=0; x<schema.max_x; x+=16){
			for (let y=0; y<schema.max_y; y+=16){
				for (let z=0; z<schema.max_z; z+=16){
					const schemapart = await schemapart_dao.get_by_id_and_offset(schema.id, x, y, z);
					if (schemapart) {
						const data = serializer.deserialize(schemapart);
						mts += export_part(data, x, y, z);
					}
				}
			}
		}

		// close outer table
		mts += "}";

		res.header("Content-Type", "application/binary").send(mts);
	} catch (e) {
		res.status(500).send(e.message);
		console.error(e);
	}
});
