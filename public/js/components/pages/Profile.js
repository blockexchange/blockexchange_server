import { get_claims } from "../../service/login.js";
import { get_user } from "../../api/user.js";

import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";
import UserProfile from "../UserProfile.js";
import UserRename from "../UserRename.js";
import LoadingBlock from "../LoadingBlock.js";
import AccessTokens from "../AccessTokens.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"user-profile": UserProfile,
		"user-rename": UserRename,
		"loading-block": LoadingBlock,
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
	methods: {
		fetch_data: function() {
			return {
				user: get_user(get_claims().user_id)
			};
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8 card" style="padding: 20px;">
				<loading-block :fetch_data="fetch_data" v-slot="{ data }">
					<user-profile :user="data.user"/>
					<hr>
					<access-tokens :username="claims.username"/>
					<hr>
					<user-rename/>
				</loading-block>
			</div>
			<div class="col-md-2"></div>
		</div>
		`
};
