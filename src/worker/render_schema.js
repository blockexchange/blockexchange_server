const SchemaRenderer = require("../render/SchemaRenderer");
const { create } = require("../dao/schema_screenshot");

console.log(process.argv);

const schema_id = +process.argv[2];

SchemaRenderer.render(schema_id)
.then(png => {
	// save png
	console.log("png", schema_id, png.length);
	create(schema_id, "Preview", "image/png", png);
});
