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

		const img1 = PImage.make(100, 100);
		const ctx = img1.getContext('2d');
		ctx.fillStyle = 'rgba(0,255,0, 0.5)';
		ctx.fillRect(0,0,100,100);
		ctx.fillStyle = 'rgba(255,0,0, 0.5)';
		ctx.fillRect(10,10,80,80);

		PImage.encodePNGToStream(img1, fs.createWriteStream('image.png')).then(() => {
		    console.log("wrote out the png file to out.png");
		}).catch((e) => {
		    console.log("there was an error writing", e);
		});
	});
});
