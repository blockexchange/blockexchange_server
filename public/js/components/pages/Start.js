import Breadcrumb from "../Breadcrumb.js";
import { START } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="text-center">
			<h3>
				Blockexchange
				<small class="text-muted">start</small>
			</h3>
			<hr/>
			<router-link to="/profile" class="btn btn-primary">
				<i class="fa fa-user"></i> Profile
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/blockexchange/blockexchange_server" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
		</div>
	`
};
