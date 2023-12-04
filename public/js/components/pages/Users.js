import Breadcrumb, { START, USERS } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, USERS]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
