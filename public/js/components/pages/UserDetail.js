import Breadcrumb, { START, USERS, USER_DETAIL } from "../Breadcrumb.js";
import UserProfile from "../UserProfile.js";
import PagedContent from "../PagedContent.js";
import SchemaList from "../SchemaList.js";

import { schema_count, schema_search } from "../../api/schema.js";

export default {
	components: {
		"bread-crumb": Breadcrumb,
		"user-profile": UserProfile,
		"schema-list": SchemaList,
        "paged-content": PagedContent
	},
	props: ["username"],
	data: function() {
		return {
			breadcrumb: [START, USERS, USER_DETAIL(this.username)]
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
        }
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4 card" style="padding: 20px;">
				<user-profile :username="username"/>
			</div>
			<div class="col-md-4"></div>
		</div>
		<hr>
		<div class="row">
			<div class="col-md-12">
				<paged-content
					:fetch_entries="fetch_entries"
					:count_entries="count_entries"
					per_page="12">
					<template #body="{ list }">
						<schema-list :list="list"/>
					</template>
				</paged-content>
			</div>
		</div>
		`
};
