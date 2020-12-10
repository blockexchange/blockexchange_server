const fs = require("fs");
const { createCanvas } = require('canvas');


describe('renderer', function() {
	it('renders an image', function() {
		const width = 1200;
		const height = 600;

		const canvas = createCanvas(width, height);
		const context = canvas.getContext('2d');

		context.fillStyle = '#ff0000';
		context.fillRect(0, 0, width, height);

		const buffer = canvas.toBuffer('image/png');
		fs.writeFileSync('./image.png', buffer);
	});
});
