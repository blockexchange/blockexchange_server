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

		console.log(get_point(0,0,0), get_color(get_point(0,0,0)));
		console.log(get_point(15,15,15), get_color(get_point(15,15,15)));

		const tan30 = Math.tan(30 * Math.PI / 180);
		const sqrt3div2 = 2 / Math.sqrt(3);

		function adjust_color_component(c, value){
			return Math.min(255, Math.max(0, c+value));
		}

		function adjust_color(color, value){
			return "rgb(" + adjust_color_component(color.r, value) +
				"," + adjust_color_component(color.g, value) +
				"," + adjust_color_component(color.b, value) + ")";
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
		const x_offset = 200;
		const y_offset = 200;
		const size = 20;

		for (let y=0; y<16; y++){
			for (let x=0; x<16; x++){
				const z = 0;
				const nodeid = get_point(x,y,z);
				const color = get_color(nodeid);

				if (color) {
					drawCube(ctx, x_offset+(size*x), y_offset+(size*tan30*y), size, color);
				}
			}
		}

		PImage.encodePNGToStream(img1, fs.createWriteStream('image.png')).then(() => {
		    console.log("wrote out the png file to out.png");
		}).catch((e) => {
		    console.log("there was an error writing", e);
		});
	});
});
