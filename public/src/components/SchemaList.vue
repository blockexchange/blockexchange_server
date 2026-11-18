<script setup lang="ts">
import format_size from "../util/format_size";
import EmptyState from "./EmptyState.vue";
import type { Schema, SchemaSearchResponse } from "../api/types";

withDefaults(defineProps<{
	list: SchemaSearchResponse[];
	show_not_found?: boolean;
	show_author?: boolean;
}>(), {
	show_not_found: true,
	show_author: true
});

defineSlots<{
	empty?(): unknown;
}>();

const preview_src = (schema: Schema) => '/api/schema/' + schema.uid + '/screenshot?height=240&width=360';
const schema_link = (entry: SchemaSearchResponse) => '/schema/' + entry.username + '/' + entry.schema.name;
</script>

<template>
	<template v-if="list && list.length == 0 && show_not_found">
		<slot name="empty">
			<empty-state icon="cubes" title="No schematics found">
				<p class="mb-0">Try different keywords or remove a filter.</p>
			</empty-state>
		</slot>
	</template>
	<div class="bx-schema-grid" v-else>
		<article class="card bx-schema-card position-relative" :key="entry.schema.uid" v-for="entry in list">
			<div class="bx-thumb">
				<img :src="preview_src(entry.schema)" :alt="'Preview of ' + entry.schema.name" loading="lazy"/>
			</div>
			<div class="card-body pb-2">
				<div class="d-flex justify-content-between align-items-start gap-2">
					<h3 class="card-title">
						<router-link :to="schema_link(entry)">{{ entry.schema.name }}</router-link>
					</h3>
					<span class="small text-nowrap" v-if="entry.schema.stars" :title="entry.schema.stars + ' stars'">
						<i class="fa fa-star bx-star"></i> {{ entry.schema.stars }}
					</span>
				</div>
				<div class="bx-meta small text-muted-bx" v-if="show_author">
					by <router-link :to="'/user/' + entry.username" class="link-body-emphasis">{{ entry.username }}</router-link>
				</div>
				<p class="small text-muted-bx mt-2 mb-2 text-truncate" v-if="entry.schema.short_description">
					{{ entry.schema.short_description }}
				</p>
				<div class="bx-meta d-flex flex-wrap gap-1 mt-2" v-if="entry.tags.length || entry.collection_name">
					<router-link
						v-for="tag in entry.tags"
						:key="tag"
						:to="{path: '/search', query: {tag_name: tag}}"
						class="badge text-bg-success text-decoration-none">
						{{ tag }}
					</router-link>
					<router-link
						v-if="entry.collection_name"
						:to="'/collections/' + entry.username + '/' + entry.collection_name"
						class="badge text-bg-info text-decoration-none"
						title="Collection">
						<i class="fa fa-object-group"></i> {{ entry.collection_name }}
					</router-link>
				</div>
			</div>
			<div class="card-footer d-flex justify-content-between">
				<span title="Width × height × depth in nodes">
					<i class="fa fa-cube"></i>
					{{ entry.schema.size_x }} × {{ entry.schema.size_y }} × {{ entry.schema.size_z }}
				</span>
				<span>{{ format_size(entry.schema.total_size) }}</span>
			</div>
		</article>
	</div>
</template>
