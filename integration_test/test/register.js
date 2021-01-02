const assert = require('assert');
const axios = require('axios');

describe('register-endpoint', function() {
	it("fail on empty password", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/register", {
			name: "user1",
			password: ""
		});
		assert.equal(false, data.success);
	});
	it("fail on empty username", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/register", {
			name: "",
			password: "password"
		});
		assert.equal(false, data.success);
	});
	it("fail on duplicate user", async function(){
		const { data } = await axios.post("http://blockexchange:8080/api/register", {
			name: "temp",
			password: "password"
		});
		assert.equal(false, data.success);
	});
});
