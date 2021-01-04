const assert = require('assert');
const axios = require('axios');

describe('register-endpoint', function() {
	const username = "integration_test_" + Math.floor(Math.random() * 10000);
	const password = username + "pw";

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
		let result = await axios.post("http://blockexchange:8080/api/register", {
			name: username,
			password: password
		});
		assert.equal(true, result.data.success);
		result = await axios.post("http://blockexchange:8080/api/register", {
			name: username,
			password: password
		});
		assert.equal(false, result.data.success);
	});
});
