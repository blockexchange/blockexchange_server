import Breadcrumb, { LOGIN, START } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, LOGIN]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
