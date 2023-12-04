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
				<small class="text-muted">login</small>
			</h3>
			<hr/>
		</div>
	`
};
