const { createCanvas } = require('canvas');

const { get_by_id_and_offset } = require("../dao/schemapart");
const { get_by_id } = require("../dao/schema");

const MapblockRenderer = require("./MapblockRenderer");
const serializer = require("../util/serializer");

module.exports.render = function(schemaid){

	const canvas = createCanvas(1024, 1024);
	const ctx = canvas.getContext('2d');

	function render_part(x, y, z){
		return get_by_id_and_offset(schemaid, x, y, z)
		.then(schemapart => {

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
			MapblockRenderer.render(ctx, mapblock, 500, 900);

			return new Promise(resolve => {
				const stream = canvas.createPNGStream();
				const bufs = [];
				stream.on('data', function(d){ bufs.push(d); });
				stream.on('end', function(){
			  	resolve(Buffer.concat(bufs));
				});
			});
		});
	}

	return get_by_id(schemaid)
	.then(function(schema){
		console.log(schema);
		return render_part(16,0,0);
	});

};
