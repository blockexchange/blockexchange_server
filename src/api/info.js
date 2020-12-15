const app = require("../app");
const logger = require("../logger");

app.get('/api/info', function(req, res){
  logger.debug("GET /api/info");

	res.json({
		api_version_major: 1,
		api_version_minor: 1,
		name: process.env.BLOCKEXCHANGE_NAME || "unknown",
		owner: process.env.BLOCKEXCHANGE_OWNER || "unknown",
		oauth: {
			github_id: process.env.GITHUB_APP_ID,
			discord_id: process.env.DISCORD_APP_ID,
			mesehub_id: process.env.MESEHUB_APP_ID,
			base_url: process.env.BASE_URL
		},
    matomo: {
      url: process.env.MATOMO_URL,
      id: process.env.MATOMO_ID
    }
	});
});
