
// user mgmt
require("./api/user");
require("./api/register");
require("./api/token");
require("./api/access_token");
require("./api/oauth/github");
require("./api/oauth/discord");
require("./api/oauth/mesehub");

// stats / info
require("./api/info");
require("./api/metrics");
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

// tags
require("./api/tag");

// collections
require("./api/collection");
require("./api/collection_schema");

// static page
require("./api/static");
