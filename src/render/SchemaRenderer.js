const { createCanvas } = require('canvas');

const { get_by_id_and_offset } = require("../dao/schemapart");
const { get_by_id } = require("../dao/schema");

const MapblockRenderer = require("./MapblockRenderer");
const serializer = require("../util/serializer");

const tan30 = Math.tan(30 * Math.PI / 180);
const sqrt3div2 = 2 / Math.sqrt(3);

module.exports.render = async function(schemaid){

	const canvas = createCanvas(1024, 1024);
	const ctx = canvas.getContext('2d');

	const schema = await get_by_id(schemaid);

	const start_block_x = Math.ceil(schema.max_x / 16) - 1;
	const start_block_z = Math.ceil(schema.max_z / 16) - 1;
	const end_block_y = Math.ceil(schema.max_y / 16) - 1;

	for (let block_x=start_block_x; block_x>=0; block_x--){
		for (let block_z=start_block_z; block_z>=0; block_z--){
			for (let block_y=0; block_y<=end_block_y; block_y++){
				const x = block_x * 16;
				const y = block_y * 16;
				const z = block_z * 16;

				const schemapart = await get_by_id_and_offset(schemaid, x, y, z);

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

				const size = 10;
				const x_offset = 500+(size*x)-(size*z);
				const y_offset = 900-(size*tan30*x)-(size*tan30*z)-(size*sqrt3div2*y);

				MapblockRenderer.render(ctx, mapblock, 10, x_offset, y_offset);

			}
		}
	}

	const png = await (new Promise(resolve => {
		const stream = canvas.createPNGStream();
		const bufs = [];
		stream.on('data', function(d){ bufs.push(d); });
		stream.on('end', function(){
			resolve(Buffer.concat(bufs));
		});
	}));

	return png;
};
