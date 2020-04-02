const assert = require('assert');
const serializer = require("../src/util/serializer");
const fs = require("fs");

describe('serializer', function() {
  describe('serialize', function() {
    it('matches simple in/output', function() {

			const data = {
			  node_ids: [
					200,0,
					200,0,
					0,0,
					0,0
				],
			  param1: [
					15,15,
					15,15,
					15,15,
					15,15
				],
			  param2: [
					0,0,
					0,0,
					0,0,
					0,0
				],
			  node_mapping: {
			    "air": 0,
			    "default:dirt": 200
			  },
			  metadata: {
			    "(0,0,0)": {
			      fields: {},
			      inventories: {}
			    }
			  },
			  size: {
			    x: 2,
			    y: 2,
			    z: 2
			  }
			};

			const schemapart = serializer.serialize(data);
			const out_data = serializer.deserialize(schemapart);

			assert.equal(data.size.x, out_data.size.x);
    });

		/*
		it('matches real-life in/output', function() {
			const data_str = fs.readFileSync("./test/schemapart.json");
			const data = JSON.parse(data_str);
			const schemapart = serializer.serialize(data);
			const out_data = serializer.deserialize(schemapart);

			assert.equal(data.size.x, out_data.size.x);
		});
		*/
  });
});
