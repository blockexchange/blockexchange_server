<script setup lang="ts">
import { computed } from 'vue';

import BreadCrumb from "../Breadcrumb.vue";
import UserProfile from "../UserProfile.vue";
import PagedContent from "../PagedContent.vue";
import SchemaList from "../SchemaList.vue";
import { START, USERS, USER_DETAIL } from "../../util/breadcrumbs";
import { schema_count, schema_search } from "../../api/schema";
import type { SchemaSearchRequest } from "../../api/types";

const props = defineProps<{
	username: string;
}>();

const breadcrumb = computed(() => [START, USERS, USER_DETAIL(props.username)]);

const search_body = (limit?: number, offset?: number): SchemaSearchRequest => ({
	user_name: props.username,
	limit: limit,
	offset: offset,
	complete: true
});

const fetch_entries = (limit: number, offset: number) => schema_search(search_body(limit, offset));
const count_entries = () => schema_count(search_body());
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="row">
		<div class="col-md-4"></div>
		<div class="col-md-4 card" style="padding: 20px;">
			<user-profile :username="username" :key="username"/>
		</div>
		<div class="col-md-4"></div>
	</div>
	<hr>
	<div class="row">
		<div class="col-md-12">
			<paged-content
				:fetch_entries="fetch_entries"
				:count_entries="count_entries"
				:per_page="12">
				<template #body="{ list }">
					<schema-list :list="list"/>
				</template>
			</paged-content>
		</div>
	</div>
</template>
