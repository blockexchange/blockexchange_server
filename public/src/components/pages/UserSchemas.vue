<script setup lang="ts">
import { computed } from 'vue';

import BreadCrumb from "../Breadcrumb.vue";
import PagedContent from "../PagedContent.vue";
import SchemaList from "../SchemaList.vue";
import { START, USER_SCHEMAS } from "../../util/breadcrumbs";
import { schema_search, schema_count } from "../../api/schema";
import type { SchemaSearchRequest } from "../../api/types";

const props = defineProps<{
	username: string;
}>();

const breadcrumb = computed(() => [START, USER_SCHEMAS(props.username)]);

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
	<paged-content
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list }">
			<schema-list :list="list"/>
		</template>
	</paged-content>
</template>
