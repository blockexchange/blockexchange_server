import Breadcrumb, { START, SEARCH } from "../Breadcrumb.js";
import SchemaSearch from "../SchemaSearch.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"schema-search": SchemaSearch
	},
	data: function() {
		return {
			breadcrumb: [START, SEARCH]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<schema-search/>
	`
};
