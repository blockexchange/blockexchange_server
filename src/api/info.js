const app = require("../app");

app.get('/api/info', function(req, res){
  console.log("GET /api/info");

	res.json({
		api_version: "1.0.0",
		name: process.env.BLOCKEXCHANGE_NAME || "unknown",
		owner: process.env.BLOCKEXCHANGE_OWNER || "unknown",
	});
});
