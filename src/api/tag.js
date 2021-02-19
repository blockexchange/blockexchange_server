
const app = require("../app");
const tag_dao = require("../dao/tag");
const logger = require("../logger");

app.get('/api/tag', async function(req, res){
	logger.debug("GET /api/tag");

	const tags = await tag_dao.find_all();
	res.json(tags);
});
