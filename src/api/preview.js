const { createCanvas } = require('canvas');

const app = require("../app");
const logger = require("../logger");

const MapblockRenderer = require("../render/MapblockRenderer");
const SchemaRenderer = require("../render/SchemaRenderer");
const serializer = require("../util/serializer");
const { get_by_id_and_offset } = require("../dao/schemapart");

// preview for a single schemapart
app.get('/api/preview/schemapart/:schema_id/:offset_x/:offset_y/:offset_z', async function(req, res){
	logger.debug("GET /api/preview/schemapart/:schema_id/:offset_x/:offset_y/:offset_z", req.params);

	const schemapart = await get_by_id_and_offset(req.params.schema_id, req.params.offset_x, req.params.offset_y, req.params.offset_z);
	if (!schemapart){
		return res.status(404).end();
	}

	const data = serializer.deserialize(schemapart);

	const mapblock = {
		data: {
			node_ids: data.node_ids,
			param1: data.param1,
			param2: data.param2,
			metadata: data.metadata,
			size: data.size,
			node_mapping: data.node_mapping
		}
	};

	const canvas = createCanvas(1024, 1024);
	const ctx = canvas.getContext('2d');
	
	MapblockRenderer.render(ctx, mapblock, 20, 500, 900);

	const png = canvas.toBuffer("image/png");
	if (png){
		res.header("Content-type", "image/png")
		.send(png);
	} else {
		res.status(404).send("not found");
	}
});

// preview for the whole schema
app.get('/api/preview/schema/:schema_id', async function(req, res){
	logger.debug("GET /api/preview/schema/:schema_id", req.params);

	const png = await SchemaRenderer.render(req.params.schema_id);
	if (png){
		res.header("Content-type", "image/png")
		.send(png);
	} else {
		res.status(404).send("not found");
	}
});
