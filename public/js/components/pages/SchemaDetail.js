import Breadcrumb, { START, SCHEMA_DETAIL } from "../Breadcrumb.js";

export default {
    props: ["username", "name"],
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SCHEMA_DETAIL(this.username, this.name)]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
	`
};
