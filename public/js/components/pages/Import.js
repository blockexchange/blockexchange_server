import Breadcrumb, { START, SCHEMA_IMPORT } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SCHEMA_IMPORT]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
