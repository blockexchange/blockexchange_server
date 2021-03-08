const assert = require('assert');
const export_part = require("../src/util/export_part");
const fs = require("fs");
//const luaparse = require("luaparse");

describe('export_part', function() {
  describe('export', function() {
    it('exports the metdata correctly', function() {

			const data_str = fs.readFileSync("./test/testdata/metadata_mapblock.json");
			const mapblock = JSON.parse(data_str);

			const mts = export_part(mapblock.data, 0, 0, 0)
			//console.log(mts);
			assert.ok(mts);
			//const ast = luaparse.parse(`return {${mts}}`);
			//assert.ok(ast);
		});
	});
});
