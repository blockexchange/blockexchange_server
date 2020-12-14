const fs = require("fs");
const colormapping = JSON.parse(fs.readFileSync("test/colormapping.json"));

const tan30 = Math.tan(30 * Math.PI / 180);
const sqrt3div2 = 2 / Math.sqrt(3);

function get_color(nodename){
	return colormapping[nodename];
}

function adjust_color_component(c, value){
	return Math.min(255, Math.max(0, c+value));
}

const cache = {};
function adjust_color(color, value){
	const key = color.r + "/" + color.g + "/" + color.b + "/" + value;
	if (cache[key]){
		return cache[key];
	}
	const str = "rgb(" + adjust_color_component(color.r, value) +
		"," + adjust_color_component(color.g, value) +
		"," + adjust_color_component(color.b, value) + ")";

	cache[key] = str;
	return str;
}

function drawCube(ctx, x, y, r, color){
	// right side
	ctx.fillStyle = adjust_color(color, 0);
	ctx.beginPath();
	ctx.moveTo(r+x, (r*tan30)+y);
	ctx.lineTo(x, (r*sqrt3div2)+y);
	ctx.lineTo(x,y);
	ctx.lineTo(r+x, -(r*tan30)+y);
	ctx.closePath();
	ctx.fill();

	// left side
	ctx.fillStyle = adjust_color(color, -20);
	ctx.beginPath();
	ctx.moveTo(x, (r*sqrt3div2)+y);
	ctx.lineTo(-r+x, (r*tan30)+y);
	ctx.lineTo(-r+x, -(r*tan30)+y);
	ctx.lineTo(x,y);
	ctx.closePath();
	ctx.fill();

	// top side
	ctx.fillStyle = adjust_color(color, 20);
	ctx.beginPath();
	ctx.moveTo(-r+x, -(r*tan30)+y);
	ctx.lineTo(x, -(r*sqrt3div2)+y);
	ctx.lineTo(r+x, -(r*tan30)+y);
	ctx.lineTo(x,y);
	ctx.closePath();
	ctx.fill();
}


module.exports.render = function(ctx, mapblock, size, x_offset, y_offset){

	if (Object.keys(mapblock.data.node_mapping).length == 1 && mapblock.data.node_mapping.air){
		// only air, skip
		return;
	}

	const y_multiplier = mapblock.data.size.z;
	const x_multiplier = mapblock.data.size.y * mapblock.data.size.z;

	const max_x = mapblock.data.size.x - 1;
	const max_z = mapblock.data.size.z - 1;
	const max_y = mapblock.data.size.y - 1;

	function get_point(x,y,z){
		if (x>max_x || y>max_y || z>max_z || x<0 || y<0 || z<0)
			//out of bounds
			return;

		const index = z + (y*y_multiplier) + (x*x_multiplier);
		return mapblock.data.node_ids[index];
	}

	// reverse index
	const nodeid_mapping = [];
	Object.keys(mapblock.data.node_mapping).forEach(function(nodename){
		nodeid_mapping[mapblock.data.node_mapping[nodename]] = nodename;
	});

	function get_image_pos_x(x,y,z){
		return x_offset+(size*x)-(size*z);
	}

	function get_image_pos_y(x,y,z){
		return y_offset-(size*tan30*x)-(size*tan30*z)-(size*sqrt3div2*y);
	}

	const blocks = [];

	function probe_position(x,y,z){
		const nodeid = get_point(x,y,z);
		const color = get_color(nodeid_mapping[nodeid]);

		if (color){
			// block to draw found, mark and return
			blocks.push({
				pos: { x:x, y:y, z:z },
				color: color,
				order: y + ((max_x-x)*max_x) + ((max_z-z)+max_z)
			});
			return;
		}

		// no color found, search in the next layer
		const next_x = x+1;
		const next_z = z+1;
		const next_y = y-1;

		// check coordinate limits
		if (next_x > max_x || next_z > max_z || next_y < 0){
			// mapblock ends
			return;
		}

		// recurse
		return probe_position(next_x, next_y, next_z);
	}

	for (let y=0; y<max_y; y++){
		// right side
		for (let x=max_x; x>=1; x--)
			probe_position(x,y,0);

		// left side
		for (let z=max_z; z>=0; z--)
			probe_position(0,y,z);
	}

	// top side
	for (let z=max_z; z>=0; z--)
		for (let x=max_x; x>=0; x--)
			probe_position(x,max_y,z);

	blocks.sort(function(a,b){
		if (a.order > b.order)
			return 1;
		else if (a.order < b.order)
			return -1;
		else
			return 0;
	});

	blocks.forEach(function(block){
		const { x, y, z } = block.pos;
		drawCube(ctx, get_image_pos_x(x,y,z), get_image_pos_y(x,y,z), size, block.color);
	});
};
