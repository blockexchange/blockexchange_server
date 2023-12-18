import Breadcrumb, { START, USERS, USER_DETAIL } from "../Breadcrumb.js";
import UserProfile from "../UserProfile.js";
import LoadingBlock from "../LoadingBlock.js";

import { search_users } from "../../api/user.js";

export default {
	components: {
		"bread-crumb": Breadcrumb,
		"user-profile": UserProfile,
		"loading-block": LoadingBlock
	},
	props: ["username"],
	data: function() {
		return {
			breadcrumb: [START, USERS, USER_DETAIL(this.username)]
		};
	},
	methods: {
		fetch_data: function() {
			return {
				user: search_users({ name: this.username }).then(l => l[0])
			};
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4 card" style="padding: 20px;">
				<loading-block :fetch_data="fetch_data" v-slot="{ data }">
					<user-profile :user="data.user"/>
				</loading-block>
			</div>
			<div class="col-md-4"></div>
		</div>
		`
};
