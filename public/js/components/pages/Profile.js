import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, PROFILE]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};