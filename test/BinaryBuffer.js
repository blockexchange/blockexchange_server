const assert = require('assert');
const BinaryBuffer = require("../src/util/BinaryBuffer");

describe('BinaryBuffer', function() {
	it('supports adding strings', function() {
		const b = new BinaryBuffer();
		b.add(0x0a);
		b.add(`H`);
		b.add(`Hello world\n\0x`);
		const buf = b.toBuffer();
		assert.equal(0x0a, buf[0]);
		assert.equal("H".charCodeAt(0), buf[1]);
		assert.equal("H".charCodeAt(0), buf[2]);
	});
});
