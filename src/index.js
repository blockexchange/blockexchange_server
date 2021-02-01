const app = require("./app");
const migrate = require("./migrate");
const logger = require("./logger");

// user mgmt
require("./api/user");
require("./api/register");
require("./api/token");
require("./api/access_token");
require("./api/oauth_callback_github");
require("./api/oauth_callback_discord");
require("./api/oauth_callback_mesehub");

// stats / info
require("./api/info");
require("./api/preview");

// search
require("./api/searchschema");

// schema metadata
require("./api/schema_star");

// down / upload
require("./api/schema_get");
require("./api/schema_update");
require("./api/schema_create");
require("./api/schema_delete");
require("./api/schema_screenshot");
require("./api/schema_export");
require("./api/schemamods");
require("./api/schemapart");

// collections
require("./api/collection");
require("./api/collection_schema");

// static page
require("./api/static");

// events
require("./events/render_schema");

// discord feed
const discord_feed = require("./feed/discord");
discord_feed();

const cleanupjob = require("./jobs/schema_cleanup");

migrate().then(() => {
	logger.info("DB Migration done");
	app.listen(8080, err => {
		if (err){
			logger.error(err);
		} else {
			logger.info('Listening on http://127.0.0.1:8080');
			cleanupjob.start();
		}
	});
})
.catch(e => {
	logger.error(e);
	cleanupjob.stop();
	process.exit(-1);
});
