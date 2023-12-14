import Breadcrumb, { START, USERS, USER_DETAIL } from "../Breadcrumb.js";
import UserProfile from "../UserProfile.js";

export default {
	components: {
		"bread-crumb": Breadcrumb,
		"user-profile": UserProfile
	},
	props: ["username"],
	data: function() {
		return {
			breadcrumb: [START, USERS, USER_DETAIL(this.username)]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4 card" style="padding: 20px;">
				<user-profile :name="username"/>
			</div>
			<div class="col-md-4"></div>
		</div>
		`
};
