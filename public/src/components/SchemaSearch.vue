<script setup lang="ts">
import { computed, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';

import { schema_search, schema_count } from "../api/schema";
import type { SchemaSearchRequest } from "../api/types";
import { get_tags } from "../service/tags";
import store from "../store/schema_search";
import debounce from "../util/debounce";

import SchemaList from "./SchemaList.vue";
import PagedContent from "./PagedContent.vue";
import EmptyState from "./EmptyState.vue";

const route = useRoute();
const router = useRouter();

const tags = computed(get_tags);

function search_body(limit?: number, offset?: number): SchemaSearchRequest {
	return {
		keywords: store.keywords ? store.keywords : null,
		tag_name: store.tag_name ? store.tag_name : null,
		mod_name: store.mod_name ? store.mod_name : null,
		complete: true,
		limit: limit,
		offset: offset
	};
}

const fetch_entries = (limit: number, offset: number) => schema_search(search_body(limit, offset));
const count_entries = () => schema_count(search_body());

const has_filter = computed(() => !!(store.keywords || store.tag_name || store.mod_name));

const filter_query = () => ({
	q: store.keywords ? store.keywords : undefined,
	tag_name: store.tag_name ? store.tag_name : undefined,
	mod_name: store.mod_name ? store.mod_name : undefined
});

// apply the filters to the url, a new search starts at the first page
function search() {
	router.replace({ path: route.path, query: filter_query() });
}

const search_debounced = debounce(search, 400);

function reset() {
	store.keywords = "";
	store.tag_name = "";
	store.mod_name = "";
	search();
}

// get data from query
const query_string = (v: unknown) => typeof v == "string" ? v : "";
if (route.query.q || route.query.tag_name || route.query.mod_name) {
	// link with filters (e.g. a tag or mod badge): use exactly these filters
	store.keywords = query_string(route.query.q);
	store.tag_name = query_string(route.query.tag_name);
	store.mod_name = query_string(route.query.mod_name);
} else if (has_filter.value) {
	// restore the previous search
	router.replace({ path: route.path, query: { ...filter_query(), page: route.query.page } });
}

// keep the form in sync when a filter link is followed while on this page
const own_path = route.path;
watch(() => route.query, q => {
	if (route.path != own_path) {
		// navigating away, keep the stored search
		return;
	}
	store.keywords = query_string(q.q);
	store.tag_name = query_string(q.tag_name);
	store.mod_name = query_string(q.mod_name);
});

watch(() => [store.keywords, store.mod_name], search_debounced);
watch(() => store.tag_name, search);
</script>

<template>
	<form class="card mb-3" @submit.prevent="search" role="search">
		<div class="card-body">
			<div class="row g-2">
				<div class="col-12 col-lg-5">
					<label class="form-label small text-muted-bx" for="search-keywords">Keywords</label>
					<div class="input-group">
						<span class="input-group-text"><i class="fa fa-magnifying-glass"></i></span>
						<input id="search-keywords" type="search" class="form-control" v-model="store.keywords" placeholder="e.g. house, tower, bridge"/>
					</div>
				</div>
				<div class="col-6 col-lg-3">
					<label class="form-label small text-muted-bx" for="search-mod">Uses mod</label>
					<input id="search-mod" type="search" class="form-control" v-model="store.mod_name" placeholder="e.g. default"/>
				</div>
				<div class="col-6 col-lg-2">
					<label class="form-label small text-muted-bx" for="search-tag">Tag</label>
					<select id="search-tag" class="form-select" v-model="store.tag_name">
						<option value="">Any tag</option>
						<option v-for="tag in tags" :key="tag.uid" :value="tag.name">{{ tag.name }}</option>
					</select>
				</div>
				<div class="col-12 col-lg-2 d-flex align-items-end">
					<button class="btn btn-success w-100" type="submit">
						Search
					</button>
				</div>
			</div>
		</div>
	</form>
	<paged-content
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list, total }">
			<div class="d-flex justify-content-between align-items-center mb-3 small text-muted-bx">
				<span>
					<b class="text-body-emphasis">{{ total }}</b> {{ total == 1 ? 'schematic' : 'schematics' }}
					<template v-if="has_filter">match your filters</template>
				</span>
				<button type="button" class="btn btn-link btn-sm p-0" v-if="has_filter" @click="reset">
					<i class="fa fa-xmark"></i> Clear filters
				</button>
			</div>
			<schema-list :list="list">
				<template #empty v-if="!has_filter">
					<empty-state icon="cubes" title="No schematics yet">
						<p class="mb-0">Be the first to share a build. <router-link to="/mod">Here's how</router-link>.</p>
					</empty-state>
				</template>
			</schema-list>
		</template>
	</paged-content>
</template>
