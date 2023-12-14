import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";
import { get_claims } from "../../service/login.js";
import UserProfile from "../UserProfile.js";
import UserRename from "../UserRename.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"user-profile": UserProfile,
		"user-rename": UserRename
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
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
				<user-profile :name="claims.username"/>
				<hr>
				<user-rename/>
			</div>
			<div class="col-md-4"></div>
		</div>
		`
};
