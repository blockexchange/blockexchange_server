import Breadcrumb, { START, USER_SCHEMAS } from "../Breadcrumb.js";
import SchemaList from "../SchemaList.js";
import PagedContent from "../PagedContent.js";

import { schema_search, schema_count } from "../../api/schema.js";

export default {
    props: ["username"],
	components: {
        "bread-crumb": Breadcrumb,
        "schema-list": SchemaList,
        "paged-content": PagedContent
	},
	data: function() {
		return {
			breadcrumb: [START, USER_SCHEMAS(this.username)]
		};
	},
    methods: {
        search_body: function(limit, offset) {
            return {
                user_name: this.username,
                limit: limit,
                offset: offset,
                complete: true
            };
        },
        fetch_entries: function(limit, offset) {
            return schema_search(this.search_body(limit, offset));
        },
        count_entries: function() {
            return schema_count(this.search_body());
        },
    },
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<paged-content
            :fetch_entries="fetch_entries"
            :count_entries="count_entries"
            per_page="24">
            <template #body="{ list }">
                <schema-list :list="list"/>
            </template>
        </paged-content>
	`
};
