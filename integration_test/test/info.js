const assert = require('assert');
const axios = require('axios');

describe('info-endpoint', function() {
	it("returns data", async function(){
		const { data } = await axios.get("http://blockexchange:8080/api/info");
		assert.equal(data.api_version_major, 1);
		assert.equal(data.api_version_minor, 1);
	});
});
