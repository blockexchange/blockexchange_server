const PImage = require('pureimage');
const fs = require("fs");

describe('renderer', function() {
	it('renders an image', function() {

		const mapblock = JSON.parse(fs.readFileSync("test/schemapart_2_(0,0,0).json"));
		const colormapping = JSON.parse(fs.readFileSync("test/colormapping.json"));

		function get_point(x,y,z){
			const index = x + (y*16) + (z*256);
			return mapblock.node_ids[index];
		}

		// reverse index
		const nodeid_mapping = [];
		Object.keys(mapblock.node_mapping).forEach(function(nodename){
			nodeid_mapping[mapblock.node_mapping[nodename]] = nodename;
		});

		function get_color(nodeid){
			const nodename = nodeid_mapping[nodeid];
			return colormapping[nodename];
		}


		const tan30 = Math.tan(30 * Math.PI / 180);
		const sqrt3div2 = 2 / Math.sqrt(3);

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

		const img1 = PImage.make(1024, 1024);
		const ctx = img1.getContext('2d');
		/*
		drawCube(ctx, 50, 50, 20, { r:200, g:0, b:0 });
		drawCube(ctx, 50+20, 50+(20*tan30), 20, { r:100, g:0, b:0 });
		drawCube(ctx, 50-20, 50+(20*tan30), 20, { r:50, g:0, b:0 });
		*/

		// mapblock pos 0,0,0 offset, extends to right/left/top
		const x_offset = 500;
		const y_offset = 900;
		const size = 20;

		function get_image_pos_x(x,y,z){
			return x_offset+(size*x)-(size*z);
		}

		function get_image_pos_y(x,y,z){
			return y_offset-(size*tan30*x)-(size*tan30*z)-(size*sqrt3div2*y);
		}

		const max_x = 15;
		const max_z = 15;
		const max_y = 15;

		function draw_mapblock_pos(x,y,z){
			const nodeid = get_point(x,y,z);
			const color = get_color(nodeid);

			if (color){
				// color found, draw cube and return
				drawCube(ctx, get_image_pos_x(x,y,z), get_image_pos_y(x,y,z), size, color);
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
			return draw_mapblock_pos(next_x, next_y, next_z);
		}

		for (let y=0; y<max_y; y++){
			// right side
			for (let x=max_x; x>=0; x--)
				draw_mapblock_pos(x,y,0);

			// left side
			for (let z=max_z; z>=0; z--)
				draw_mapblock_pos(0,y,z);
		}

		// top side
		for (let z=max_z; z>=0; z--)
			for (let x=max_x; x>=0; x--)
				draw_mapblock_pos(x,max_y,z);


		PImage.encodePNGToStream(img1, fs.createWriteStream('image.png')).then(() => {
		    console.log("wrote out the png file to out.png");
		}).catch((e) => {
		    console.log("there was an error writing", e);
		});
	});
});
