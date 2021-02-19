const app = require("../app");
const registry = require("../registry");

app.get('/metrics', async function(req, res){
	res
	.header("Content-Type", "text/plain")
	.send(await registry.metrics());
});
