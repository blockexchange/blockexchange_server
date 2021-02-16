
module.exports = function export_part(data, offset_x, offset_y, offset_z){
	let mts = "";

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
					mts += `{["x"]=${offset_x+x},["y"]=${offset_y+y},["z"]=${offset_z+z},`;
					mts += `["name"]="${nodename}",["param1"]=${data.param1[index]}`;
					if (data.param2[index] > 0){
						mts += `,["param2"]=${data.param2[index]}`;
					}
					//metadata
					if (!data.metadata || !data.metadata.meta){
						mts += `},`;
						continue;
					}
					const pos_str = `(${x},${y},${z})`;
					const meta = data.metadata.meta[pos_str];
					if (meta) {
						//TODO: handle "delimiter": \u001b(T@default)\"\u001bFtest 123\u001bE\"\u001bE
						mts += `,["meta"]={`;
						if (meta.fields){
							mts += `["fields"]={`;
							mts += Object.keys(meta.fields).reduce((all, key) => {
								all += `["${key}"]="${meta.fields[key]}",`;
							});
							mts += `},`;
						}
						if (meta.inventory){
							mts += `["inventory"]={`;
							mts += Object.keys(meta.inventory).reduce((all, inv_name) => {
								all += `["${inv_name}"]={`;
								const inv_size = Object.keys(meta.inventory[inv_name]).length;
								for (let j=0; j<inv_size; j++){
									all += `"${meta.inventory[inv_name][j]}",`;
								}
								all += `},`;
							});
							mts += `}`;
						}
						mts += `}`;
					}
					mts += `},`;
				}

				index++;
			}
		}
	}

	return mts;
};
