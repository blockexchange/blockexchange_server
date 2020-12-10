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

		function drawCube(ctx, x, y, r, color){
			// right side
			ctx.fillStyle = "rgb(" + color.r + "," + color.g + "," + color.b + ")";
			ctx.beginPath();
			ctx.moveTo(r+x, (r*tan30)+y);
			ctx.lineTo(x, (r*sqrt3div2)+y);
			ctx.lineTo(x,y);
			ctx.lineTo(r+x, -(r*tan30)+y);
			ctx.closePath();
			ctx.fill();

			// left side
			ctx.fillStyle = "rgb(" + (color.r-20) + "," + color.g + "," + color.b + ")";
			ctx.beginPath();
			ctx.moveTo(x, (r*sqrt3div2)+y);
			ctx.lineTo(-r+x, (r*tan30)+y);
			ctx.lineTo(-r+x, -(r*tan30)+y);
			ctx.lineTo(x,y);
			ctx.closePath();
			ctx.fill();

			// top side
			ctx.fillStyle = "rgb(" + (color.r+20) + "," + color.g + "," + color.b + ")";
			ctx.beginPath();
			ctx.moveTo(-r+x, -(r*tan30)+y);
			ctx.lineTo(x, -(r*sqrt3div2)+y);
			ctx.lineTo(r+x, -(r*tan30)+y);
			ctx.lineTo(x,y);
			ctx.closePath();
			ctx.fill();
		}

		const img1 = PImage.make(100, 100);
		const ctx = img1.getContext('2d');
		drawCube(ctx, 50, 50, 20, { r:200, g:0, b:0 });
		drawCube(ctx, 50+20, 50+(20*tan30), 20, { r:100, g:0, b:0 });
		drawCube(ctx, 50-20, 50+(20*tan30), 20, { r:50, g:0, b:0 });

		PImage.encodePNGToStream(img1, fs.createWriteStream('image.png')).then(() => {
		    console.log("wrote out the png file to out.png");
		}).catch((e) => {
		    console.log("there was an error writing", e);
		});
	});
});
