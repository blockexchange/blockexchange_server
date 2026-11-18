<script setup lang="ts">
import { computed } from 'vue';

import PageHeader from "../PageHeader.vue";
import PagedContent from "../PagedContent.vue";
import SchemaList from "../SchemaList.vue";
import EmptyState from "../EmptyState.vue";
import { START, USER_COLLECTION, USER_COLLECTIONS } from "../../util/breadcrumbs";
import { schema_search, schema_count } from "../../api/schema";
import type { SchemaSearchRequest } from "../../api/types";

const props = defineProps<{
	username: string;
	collection_name: string;
}>();

const breadcrumb = computed(() => [START, USER_COLLECTIONS(props.username), USER_COLLECTION(props.username, props.collection_name)]);

const search_body = (limit?: number, offset?: number): SchemaSearchRequest => ({
	collection_name: props.collection_name,
	user_name: props.username,
	limit: limit,
	offset: offset,
	complete: true
});

const fetch_entries = (limit: number, offset: number) => schema_search(search_body(limit, offset));
const count_entries = () => schema_count(search_body());
</script>

<template>
	<page-header :title="collection_name" :breadcrumb="breadcrumb">
		<span class="text-muted-bx">
			A collection by <router-link :to="'/user/' + username">{{ username }}</router-link>
		</span>
	</page-header>
	<paged-content
		:key="username + '/' + collection_name"
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list }">
			<schema-list :list="list" :show_author="false">
				<template #empty>
					<empty-state icon="object-group" title="This collection is empty">
						Add schematics to it from their edit view.
					</empty-state>
				</template>
			</schema-list>
		</template>
	</paged-content>
</template>
