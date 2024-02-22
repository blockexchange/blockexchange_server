import Breadcrumb, { START, USER_SCHEMAS, SCHEMA_DETAIL } from "../Breadcrumb.js";
import LoadingBlock from "../LoadingBlock.js";
import SchemaDetail from "../schemadetail/SchemaDetail.js";

import { has_permission, get_user_uid } from "../../service/login.js";
import { schema_search } from "../../api/schema.js";

export default {
    props: ["username", "name"],
	components: {
        "bread-crumb": Breadcrumb,
		"loading-block": LoadingBlock,
		"schema-detail": SchemaDetail
	},
	data: function() {
		return {
			breadcrumb: [START, USER_SCHEMAS(this.username), SCHEMA_DETAIL(this.username, this.name)],
			BaseURL
		};
	},
	methods: {
		fetch_data: function() {
			return {
				search_result: schema_search({
					schema_name: this.name,
					user_name: this.username
				}).then(r => r[0])
			};
		},
		allow_edit: function(schema) {
			return (get_user_uid() == schema.user_uid || has_permission("ADMIN"));
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<loading-block :fetch_data="fetch_data" v-slot="{ data, update_data }">
			<schema-detail
				:search_result="data.search_result"
				:allow_edit="allow_edit(data.search_result.schema)"
				v-on:save="update_data"
				/>
		</loading-block>
	`
};
