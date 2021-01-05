const assert = require('assert');
const axios = require('axios');
const fs = require("fs");

describe('upload-usecase', function() {
	const username = "integration_test_" + Math.floor(Math.random() * 10000);
	const password = username + "pw";
	let token = null;
	let schema = null;

	it("can register", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/register", {
			name: username,
			password: password
		});
		assert.equal(true, data.success);
	});

	it("can log in", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/token", {
			name: username,
			password: password
		});
		token = data;
		assert.ok(token);
	});

	it("can create a schema", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/schema", {
			name: "my_schema",
			description: "some long description",
			max_x: 16,
			max_y: 16,
			max_z: 16,
			part_length:16,
			license: "CC0"
		}, {
			headers: {
				"Authorization": token
			}
		});
		schema = data;
		assert.ok(schema);
		assert.ok(schema.id);
	});

	it("can upload schemaparts", async function(){
		const data_str = fs.readFileSync("./testdata/testcube.json");
		const mapblock = JSON.parse(data_str);
		mapblock.schema_id = schema.id;

		const { data } = await axios.post("http://blockexchange:8080/api/schemapart", mapblock, {
			headers: {
				"Authorization": token
			}
		});

		assert.ok(data);
		assert.ok(data.id);
	});

	it("can set used mods", async function(){
		const response = await axios.post(`http://blockexchange:8080/api/schema/${schema.id}/mods`, {
			"default": 32,
			"air": 666
		}, {
			headers: {
				"Authorization": token
			}
		});

		assert.equal(200, response.status);
	});

	it("can finalize the schema", async function(){
		const response = await axios.post(`http://blockexchange:8080/api/schema/${schema.id}/complete`, {}, {
			headers: {
				"Authorization": token
			}
		});

		assert.equal(200, response.status);
	});

	// replace

	it("can create another schema", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/schema", {
			name: "my_schema_2",
			description: "some long description",
			max_x: 16,
			max_y: 16,
			max_z: 16,
			part_length:16,
			license: "CC0"
		}, {
			headers: {
				"Authorization": token
			}
		});
		schema = data;
		assert.ok(schema);
		assert.ok(schema.id);
	});

	it("can upload schemaparts", async function(){
		const data_str = fs.readFileSync("./testdata/testcube.json");
		const mapblock = JSON.parse(data_str);
		mapblock.schema_id = schema.id;

		const { data } = await axios.post("http://blockexchange:8080/api/schemapart", mapblock, {
			headers: {
				"Authorization": token
			}
		});

		assert.ok(data);
		assert.ok(data.id);
	});

	it("can set used mods", async function(){
		const response = await axios.post(`http://blockexchange:8080/api/schema/${schema.id}/mods`, {
			"default": 32,
			"air": 666
		}, {
			headers: {
				"Authorization": token
			}
		});

		assert.equal(200, response.status);
	});

	it("can finalize the schema", async function(){
		const response = await axios.post(`http://blockexchange:8080/api/schema/${schema.id}/complete`, {
			replaces: "my_schema"
		}, {
			headers: {
				"Authorization": token
			}
		});

		assert.equal(200, response.status);
	});

});
