const app = require("../app");
const logger = require("../logger");
const { get_by_schemaname_and_username } = require("../dao/schema");
const { find_all } = require("../dao/schema_screenshot");

app.get('/api/static/schema/:username/:schemaname', async function(req, res){
  logger.debug("GET /api/static/schema", req.params);

	const schema = await get_by_schemaname_and_username(req.params.schemaname, req.params.username);

	if (!schema){
		res.status(404).end();
		return;
	}

	const screenshots = await find_all(schema.id);

	if (!screenshots || screenshots.length == 0){
		res.status(404).end();
		return;
	}


	const html = /*html*/`
	<!DOCTYPE HTML>
	<html>
		<head>
			<meta name="og:title" content="${schema.name} by ${req.params.username}"/>
			<meta name="og:type" content="Schematic"/>
			<meta name="og:url" content="${process.env.BASE_URL}/#/schema/${req.params.username}/${schema.name}"/>
			<meta name="og:image" content="${process.env.BASE_URL}/api/schema/${schema.id}/screenshot/${screenshots[0].id}"/>
			<meta name="og:site_name" content="Block exchange"/>
			<meta name="og:description" content="${schema.description}"/>
			<meta http-equiv="refresh" content="0; url=${process.env.BASE_URL}/#/schema/${req.params.username}/${schema.name}" />
		</head>
		<body>
		</body>
	</html>
	`;

	res.header("Content-type", "text/html").send(html);
});
