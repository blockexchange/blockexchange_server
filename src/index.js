const app = require("./app");

require("./api/schema");
require("./api/schemapart");

app.listen(8080, () => console.log('Listening on http://127.0.0.1:8080'));
