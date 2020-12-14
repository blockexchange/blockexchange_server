const { createCanvas } = require('canvas');

const app = require("../app");
const logger = require("../logger");

const MapblockRenderer = require("../render/MapblockRenderer");
const serializer = require("../util/serializer");
const { get_by_id_and_offset } = require("../dao/schemapart");

app.get('/api/preview/schemapart/:schema_id/:offset_x/:offset_y/:offset_z', function(req, res){
	logger.debug("GET /api/preview/schemapart/:schema_id/:offset_x/:offset_y/:offset_z", req.params);

	get_by_id_and_offset(req.params.schema_id, req.params.offset_x, req.params.offset_y, req.params.offset_z)
	.then(schemapart => {
		if (!schemapart){
			return;
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

		MapblockRenderer.render(ctx, mapblock, 500, 900);

		return new Promise(resolve => {
			const stream = canvas.createPNGStream();
			const bufs = [];
			stream.on('data', function(d){ bufs.push(d); });
			stream.on('end', function(){
				resolve(Buffer.concat(bufs));
			});
		});
	})
	.then(png => {
		if (png){
			res.header("Content-type", "image/png")
			.send(png);
		} else {
			res.status(404).send("not found");
		}
	})
	.catch(e => {
		res.status(500).send(e.message);
		console.error(e);
	});


});
