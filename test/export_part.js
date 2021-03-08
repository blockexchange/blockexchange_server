const assert = require('assert');
const export_part = require("../src/util/export_part");
const fs = require("fs");
//const luaparse = require("luaparse");

describe('export_part', function() {
  describe('export', function() {
    it('exports the metdata correctly', function() {

			const data_str = fs.readFileSync("./test/testdata/metadata_mapblock.json");
			const mapblock = JSON.parse(data_str);

			const buf = export_part(mapblock.data, 0, 0, 0);
			assert.ok(buf);
			console.log(buf);
			const stmt = `return {${buf.toString()}}`;
			console.log(stmt);
			//const ast = luaparse.parse(stmt);
			//assert.ok(ast);
		});
	});
});
