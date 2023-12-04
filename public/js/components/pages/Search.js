import Breadcrumb, { START, SEARCH } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SEARCH]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
