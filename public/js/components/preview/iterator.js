
export default function(schema){

	var pos;

	const x_parts = Math.ceil(schema.size_x / schema.part_length);
	const y_parts = Math.ceil(schema.size_y / schema.part_length);
	const z_parts = Math.ceil(schema.size_z / schema.part_length);

	return function(){
		if (!pos){
			pos = { x: 0, y: 0, z: 0 };

		} else {
			pos.x++;
			if (pos.x >= x_parts){
				pos.x = 0;
				pos.z++;
				if (pos.z >= z_parts){
					pos.z = 0;
					pos.y++;
					if (pos.y >= y_parts){
						pos = null;
					}
				}
			}
		}

		return pos;
	};
}
