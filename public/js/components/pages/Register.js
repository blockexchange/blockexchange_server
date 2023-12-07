import Breadcrumb, { START, MOD } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, MOD]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
