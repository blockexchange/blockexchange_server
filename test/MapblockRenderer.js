const assert = require('assert');
const fs = require('fs');
const { createCanvas } = require('canvas');

const MapblockRenderer = require("../src/render/MapblockRenderer");

describe('mapblockrenderer', function() {
  it('renders a mapblock', function() {
		const img_size_x = 1280;
		const img_size_y = 1024;
		const size = 10;

		const data_str = fs.readFileSync("./test/random_mapblock.json");
		const mapblock = JSON.parse(data_str);

		const canvas = createCanvas(img_size_x, img_size_y);
		const ctx = canvas.getContext('2d');

		MapblockRenderer.render(ctx, mapblock, size, img_size_x/2, img_size_y/2);

		const buf = canvas.toBuffer("image/png");
		assert.equal(true, buf != null);
  });
});
