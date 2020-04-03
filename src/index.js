const app = require("./app");

// user mgmt
require("./api/register");
require("./api/token");

// stats / info
require("./api/info");

// down / upload
require("./api/schema");
require("./api/schemapart");

app.listen(8080, () => console.log('Listening on http://127.0.0.1:8080'));
