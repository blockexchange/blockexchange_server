<script setup lang="ts">
import { computed } from 'vue';

import PageHeader from "../PageHeader.vue";
import PagedContent from "../PagedContent.vue";
import SchemaList from "../SchemaList.vue";
import EmptyState from "../EmptyState.vue";
import { START, STARRED } from "../../util/breadcrumbs";
import { starred_schema_search, starred_schema_count } from "../../api/schema_star";
import type { SchemaSearchRequest } from "../../api/types";
import { is_logged_in } from "../../service/login";

const breadcrumb = [START, STARRED];
const logged_in = computed(is_logged_in);

const search_body = (limit?: number, offset?: number): SchemaSearchRequest => ({
	limit: limit,
	offset: offset,
	complete: true
});

const fetch_entries = (limit: number, offset: number) => starred_schema_search(search_body(limit, offset));
const count_entries = () => starred_schema_count(search_body());
</script>

<template>
	<page-header
		title="Starred schematics"
		subtitle="Schematics you have starred. Only you can see this list."
		:breadcrumb="breadcrumb"/>

	<empty-state icon="lock" title="Log in to see your starred schematics" v-if="!logged_in">
		<router-link to="/login" class="btn btn-success mt-2">Log in</router-link>
	</empty-state>

	<paged-content
		v-else
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list }">
			<schema-list :list="list" :show_author="true">
				<template #empty>
					<empty-state icon="star" title="No starred schematics yet">
						<p>Star a schematic on its detail page to find it here again.</p>
						<router-link to="/search">Browse schematics</router-link>
					</empty-state>
				</template>
			</schema-list>
		</template>
	</paged-content>
</template>
