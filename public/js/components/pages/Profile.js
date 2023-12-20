import { get_claims } from "../../service/login.js";

import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";
import UserProfile from "../UserProfile.js";
import UserRename from "../UserRename.js";
import AccessTokens from "../AccessTokens.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"user-profile": UserProfile,
		"user-rename": UserRename,
		"access-tokens": AccessTokens
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
		<div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8 card" style="padding: 20px;">
				<user-profile :username="claims.username"/>
				<hr>
				<access-tokens :username="claims.username"/>
				<hr>
				<user-rename/>
			</div>
			<div class="col-md-2"></div>
		</div>
		`
};
