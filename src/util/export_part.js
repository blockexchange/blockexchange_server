
function export_node(nodename, data, index, x, y, z){
	const buf = [];

	buf.push(`{["x"]=${x},["y"]=${y},["z"]=${z},`.split(""));
	buf.push(`["name"]="${nodename}"`.split(""));
	if (data.param1[index] > 0){
		buf.push(`,["param1"]=${data.param1[index]}`.split(""));
	}
	if (data.param2[index] > 0){
		buf.push(`,["param2"]=${data.param2[index]}`.split(""));
	}
	//metadata
	if (!data.metadata || !data.metadata.meta){
		buf.push(`},`.split(""));
		return buf;
	}
	const pos_str = `(${x},${y},${z})`;
	const meta = data.metadata.meta[pos_str];
	if (meta) {
		//TODO: handle "delimiter": \u001b(T@default)\"\u001bFtest 123\u001bE\"\u001bE
		buf.push(`,["meta"]={`.split(""));
		if (meta.fields){
			buf.push(`["fields"]={`.split(""));
			Object.keys(meta.fields).forEach(key => {
				buf.push(`["${key}"]="`.split(""));
				buf.push(meta.fields[key]);
				buf.push(`"`.split(""));
			});
			buf.push(`},`.split(""));
		}
		if (meta.inventory){
			buf.push(`["inventory"]={`.split(""));
			Object.keys(meta.inventory).forEach(inv_name => {
				buf.push(`["${inv_name}"]={`.split(""));
				const inv_size = Object.keys(meta.inventory[inv_name]).length;
				for (let j=0; j<inv_size; j++){
					buf.push(`"${meta.inventory[inv_name][j]}",`.split(""));
				}
				buf.push(`},`.split(""));
			});
			buf.push(`}`.split(""));
		}
		buf.push(`}`.split(""));
	}
	buf.push(`},\n`.split(""));

	return buf;
}

module.exports = function export_part(data, offset_x, offset_y, offset_z){
	const buf = [];

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

					const parts = export_node(
						nodename,
						data,
						index,
						offset_x+x,
						offset_y+y,
						offset_z+z
					);

					parts.forEach(p => buf.push(p));
				}
				index++;
			}
		}
	}

	const all_parts = [].concat(...buf);
	console.log(all_parts);
	return Buffer.from(all_parts);
};
