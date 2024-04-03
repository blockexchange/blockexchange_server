import { get_username, get_user_type, get_user_uid } from "../../service/login.js";

import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";
import UserProfile from "../UserProfile.js";
import UserRename from "../UserRename.js";
import AccessTokens from "../AccessTokens.js";
import ChangePassword from "../ChangePassword.js";
import UnlinkOauth from "../UnlinkOauth.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"user-profile": UserProfile,
		"user-rename": UserRename,
		"access-tokens": AccessTokens,
		"change-password": ChangePassword,
		"unlink-oauth": UnlinkOauth
	},
	data: function() {
		return {
			breadcrumb: [START, PROFILE]
		};
	},
	computed: {
		useruid: get_user_uid,
		username: get_username,
		usertype: get_user_type
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-2"></div>
			<div class="col-md-8 card" style="padding: 20px;">
				<user-profile :username="username"/>
				<hr>
				<access-tokens :username="username"/>
				<hr>
				<change-password v-if="usertype == 'LOCAL'" :useruid="useruid"/>
				<unlink-oauth v-if="usertype != 'LOCAL'" :useruid="useruid"/>
				<user-rename/>
			</div>
			<div class="col-md-2"></div>
		</div>
		`
};
