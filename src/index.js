const app = require("./app");
const migrate = require("./migrate");
const logger = require("./logger");

// load api modules
require("./api");

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
