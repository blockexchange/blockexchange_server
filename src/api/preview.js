const app = require("../app");
const logger = require("../logger");
const SchemaRenderer = require("../render/SchemaRenderer");

app.get('/api/preview/:schemaid', function(req, res){
	logger.debug("GET /api/preview/:schemaid", req.params.schemaid);


	SchemaRenderer.render(+req.params.schemaid)
	.then(png => {
		res.header("Content-type", "image/png")
		.send(png);
	});
});
