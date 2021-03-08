const BinaryBuffer = require("./BinaryBuffer");

function export_node(buf, nodename, data, index, x, y, z){

	buf.add(`{["x"]=${x},["y"]=${y},["z"]=${z},`);
	buf.add(`["name"]="${nodename}"`);
	if (data.param1[index] > 0){
		buf.add(`,["param1"]=${data.param1[index]}`);
	}
	if (data.param2[index] > 0){
		buf.add(`,["param2"]=${data.param2[index]}`);
	}
	//metadata
	if (!data.metadata || !data.metadata.meta){
		buf.add("},");
		return;
	}
	const pos_str = `(${x},${y},${z})`;
	const meta = data.metadata.meta[pos_str];
	if (meta) {
		//TODO: handle "delimiter": \u001b(T@default)\"\u001bFtest 123\u001bE\"\u001bE
		buf.add(`,["meta"]={`);
		if (meta.fields){
			buf.add(`["fields"]={`);
			Object.keys(meta.fields).forEach(key => {
				buf.add(`["${key}"]="`);
				buf.add(meta.fields[key]);
				buf.add(`",`);
			});
			buf.add(`},`);
		}
		if (meta.inventory){
			buf.add(`["inventory"]={`);
			Object.keys(meta.inventory).forEach(inv_name => {
				buf.add(`["${inv_name}"]={`);
				const inv_size = Object.keys(meta.inventory[inv_name]).length;
				for (let j=0; j<inv_size; j++){
					buf.add(`"`);
					buf.add(meta.inventory[inv_name][j]);
					buf.add(`",`);
				}
				buf.add(`},`);
			});
			buf.add(`}`);
		}
		buf.add(`}`);
	}
	buf.add(`},`);
}

module.exports = function export_part(data, offset_x, offset_y, offset_z){
	const buf = new BinaryBuffer();

	// reverse node-id mapping for lookup
	const nodeid_to_name_mapping = {};
	Object.keys(data.node_mapping)
		.map(name => nodeid_to_name_mapping[data.node_mapping[name]] = name);

	const air_nodeid = data.node_mapping.air;
	const data_size = data.size.x * data.size.y * data.size.z;
	if (data_size != data.node_ids.length){
		throw new Error("unexpected data-size: " + data.node_ids.length);
	}

	let index = 0;
	for (let x=0; x<data.size.x; x++){
		for (let y=0; y<data.size.y; y++){
			for (let z=0; z<data.size.z; z++){
				const nodeid = data.node_ids[index];
				if (nodeid != air_nodeid){
					// not an air node, export
					const nodename = nodeid_to_name_mapping[nodeid];

					export_node(
						buf,
						nodename,
						data,
						index,
						offset_x+x,
						offset_y+y,
						offset_z+z
					);
				}
				index++;
			}
		}
	}

	return buf.toBuffer();
};
