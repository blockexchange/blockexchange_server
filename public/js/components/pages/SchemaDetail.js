import Breadcrumb, { START, SCHEMA_DETAIL } from "../Breadcrumb.js";
import LoadingBlock from "../LoadingBlock.js";

import { get_schema_by_name } from "../../api/schema.js";

export default {
    props: ["username", "name"],
	components: {
        "bread-crumb": Breadcrumb,
		"loading-block": LoadingBlock
	},
	data: function() {
		return {
			breadcrumb: [START, SCHEMA_DETAIL(this.username, this.name)]
		};
	},
	methods: {
		fetch_data: function() {
			return {
				schema: get_schema_by_name(this.username, this.name)
			};
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<loading-block :fetch_data="fetch_data" v-slot="{ data }">
			{{data.schema.description}}
		</loading-block>
	`
};
