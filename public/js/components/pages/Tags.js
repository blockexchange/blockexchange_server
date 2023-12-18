import Breadcrumb, { START, TAGS } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, TAGS]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
