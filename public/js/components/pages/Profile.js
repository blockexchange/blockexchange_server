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
				<small class="text-muted">profile</small>
			</h3>
			<hr/>
		</div>
	</default-layout>
	`
};
