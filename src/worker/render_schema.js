const SchemaRenderer = require("../render/SchemaRenderer");
const { create } = require("../dao/schema_screenshot");

const schema_id = +process.argv[2];

SchemaRenderer.render(schema_id)
.then(png => {
	// save png
	create(schema_id, "ISO_X+Z+", "image/png", png);
});
