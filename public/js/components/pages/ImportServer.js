
import Breadcrumb, { START, SERVER_IMPORT } from "../Breadcrumb.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	data: function() {
		return {
			breadcrumb: [START, SERVER_IMPORT]
		};
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
        <div class="container">
            <h5>Import schematic from a server</h5>
        </div>
		`
};
