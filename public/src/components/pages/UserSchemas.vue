<script setup lang="ts">
import { computed } from 'vue';

import PageHeader from "../PageHeader.vue";
import PagedContent from "../PagedContent.vue";
import SchemaList from "../SchemaList.vue";
import { START, USERS, USER_DETAIL, USER_SCHEMAS } from "../../util/breadcrumbs";
import { schema_search, schema_count } from "../../api/schema";
import type { SchemaSearchRequest } from "../../api/types";

const props = defineProps<{
	username: string;
}>();

const breadcrumb = computed(() => [START, USERS, USER_DETAIL(props.username), USER_SCHEMAS(props.username)]);

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
	<page-header :title="'Schematics by ' + username" :breadcrumb="breadcrumb">
		<template #actions>
			<router-link :to="'/user/' + username" class="btn btn-outline-light">
				<i class="fa fa-user"></i> View profile
			</router-link>
		</template>
	</page-header>
	<paged-content
		:key="username"
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list }">
			<schema-list :list="list" :show_author="false"/>
		</template>
	</paged-content>
</template>
