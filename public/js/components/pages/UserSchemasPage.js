import UserSchemas from "../search/UserSchemas.js";
import PageTitle from "../PageTitle.js";

export default {
	components: {
		"user-schemas": UserSchemas,
		"page-title": PageTitle
	},
	template: /*html*/`
		<div>
			<page-title major="User-schemas" :minor="'from ' + $route.params.user_name"/>
			<user-schemas :user_name="$route.params.user_name"/>
		</div>
	`
};
