import Breadcrumb, { START, PROFILE } from "../Breadcrumb.js";
import { get_claims } from "../../service/login.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
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
				<h5>
					User profile for
					<small class="text-body-secondary">{{claims.username}}</small>
				</h5>
				<ul>
					<li>
						<b>Username:</b> {{claims.username}}
					</li>
					<li>
						<b>ID:</b> <span class="badge bg-success">{{claims.user_id}}</span>
					</li>
					<li>
						<b>Type:</b> <span class="badge bg-secondary">{{claims.type}}</span>
					</li>
					<li>
						<b>Permissions:</b>
						<span class="badge bg-secondary" v-for="permission in claims.permissions">{{permission}}</span>
					</li>
				</ul>
            </div>
            <div class="col-md-4"></div>
		</div>
	`
};
