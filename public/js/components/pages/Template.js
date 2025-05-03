import Breadcrumb, { START } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
