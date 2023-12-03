import DefaultLayout from "../layouts/DefaultLayout.js";
import { START } from "../Breadcrumb.js";

export default {
	components: {
		"default-layout": DefaultLayout
	},
	data: function() {
		return {
			breadcrumb: [START]
		};
	},
	template: /*html*/`
	<default-layout icon="home" title="Start" :breadcrumb="breadcrumb">
		<div class="text-center">
			<h3>
				Blockexchange
				<small class="text-muted">start</small>
			</h3>
			<hr/>
			<router-link to="/" class="btn btn-primary">
				<i class="fa fa-user"></i> Start
			</router-link>
			&nbsp;
			<a class="btn btn-secondary" href="https://github.com/blockexchange/blockexchange_server" target="new">
				<i class="fa-brands fa-github"></i> Source
			</a>
		</div>
	</default-layout>
	`
};
