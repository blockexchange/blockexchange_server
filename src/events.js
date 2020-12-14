
/*
Events:

# "new-schema": <schema-object>
Created on schema finalization

# "preview-rendered": <schema-id>
Preview for the schema created

*/

const EventEmitter = require("events");
module.exports = new EventEmitter();
