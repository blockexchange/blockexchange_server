import Breadcrumb, { START, USER_SCHEMAS, SCHEMA_DETAIL } from "../Breadcrumb.js";
import LoadingBlock from "../LoadingBlock.js";
import SchemaDetail from "../schemadetail/SchemaDetail.js";

import { has_permission, get_user_uid } from "../../service/login.js";
import { get_schema_by_name } from "../../api/schema.js";

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
				search_result: get_schema_by_name(this.username, this.name)
			};
		},
		allow_edit: function(schema) {
			return (get_user_uid() == schema.user_uid || has_permission("ADMIN"));
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<loading-block :fetch_data="fetch_data" v-slot="{ data }">
			<schema-detail
				:search_result="data.search_result"
				:allow_edit="allow_edit(data.search_result.schema)"/>
		</loading-block>
	`
};
