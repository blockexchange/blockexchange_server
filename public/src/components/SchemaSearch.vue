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

const get_count = debounce(() => {
	store.busy = true;
	schema_count(search_body())
	.then(c => store.count = c)
	.then(() => store.busy = false);
}, 250);

function search() {
	// set data on query
	router.push({
		path: route.path,
		query: {
			q: store.keywords ? store.keywords : undefined,
			tag_name: store.tag_name ? store.tag_name : undefined,
			mod_name: store.mod_name ? store.mod_name : undefined,
			page: route.query.page
		}
	});
}

// get data from query
const query_string = (v: unknown) => typeof v == "string" ? v : "";
if (route.query.q) {
	store.keywords = query_string(route.query.q);
}
if (route.query.tag_name) {
	store.tag_name = query_string(route.query.tag_name);
}
if (route.query.mod_name) {
	store.mod_name = query_string(route.query.mod_name);
}

watch(() => [store.keywords, store.mod_name, store.tag_name], get_count);

// count and search on route load
get_count();
search();
</script>

<template>
	<div class="row">
		<div class="col-md-6 col-xs-4">
			<input type="text" class="form-control" v-model="store.keywords" @keyup.enter="search" placeholder="Keywords"/>
		</div>
		<div class="col-md-2 col-xs-4">
			<input type="text" class="form-control" v-model="store.mod_name" @keyup.enter="search" placeholder="Modname"/>
		</div>
		<div class="col-md-2 col-xs-4">
			<select class="form-control" v-model="store.tag_name">
				<option value="">All tags</option>
				<option v-for="tag in tags" :key="tag.uid" :value="tag.name">{{ tag.name }}</option>
			</select>
		</div>
		<div class="col-md-2 col-xs-4">
			<button class="btn btn-success w-100" @click="search">
				<i class="fa fa-spinner fa-spin" v-if="store.busy"></i>
				<i class="fa-solid fa-magnifying-glass" v-else></i>
				Search
				<span class="badge rounded-pill bg-secondary">{{ store.count }}</span>
			</button>
		</div>
	</div>
	<hr>
	<paged-content
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list }">
			<schema-list :list="list"/>
		</template>
	</paged-content>
</template>
