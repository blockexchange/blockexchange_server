const app = require("./app");
const migrate = require("./migrate");
const logger = require("./logger");

// user mgmt
require("./api/user");
require("./api/register");
require("./api/token");
require("./api/oauth_callback");

// stats / info
require("./api/info");

// search
require("./api/searchschema");

// schema metadata
require("./api/schema_star");

// down / upload
require("./api/schema_get");
require("./api/schema_update");
require("./api/schema_create");
require("./api/schema_delete");
require("./api/schemamods");
require("./api/schemapart");

// discord feed
const discord_feed = require("./feed/discord");
discord_feed();

const cleanupjob = require("./jobs/temp_schema_cleanup");

migrate().then(() => {
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
