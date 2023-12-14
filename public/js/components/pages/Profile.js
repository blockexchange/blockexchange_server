import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";
import { get_claims } from "../../service/login.js";
import UserProfile from "../UserProfile.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"user-profile": UserProfile
	},
	data: function() {
		return {
			breadcrumb: [START, PROFILE]
		};
	},
	computed: {
		claims: get_claims
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<user-profile :name="claims.username"/>
	`
};
