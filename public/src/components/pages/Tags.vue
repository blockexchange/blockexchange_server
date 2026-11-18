<script setup lang="ts">
import { computed } from 'vue';
import PageHeader from "../PageHeader.vue";
import EmptyState from "../EmptyState.vue";
import { START, TAGS } from "../../util/breadcrumbs";
import { get_tags } from "../../service/tags";

const breadcrumb = [START, TAGS];
const tags = computed(get_tags);
</script>

<template>
	<page-header title="Tags" subtitle="Tags help people find schematics. Restricted tags can only be set by admins." :breadcrumb="breadcrumb"/>
	<div class="list-group" v-if="tags.length">
		<router-link
			v-for="tag in tags"
			:key="tag.uid"
			:to="{path: '/search', query: {tag_name: tag.name}}"
			class="list-group-item list-group-item-action">
			<div class="d-flex justify-content-between align-items-center">
				<span><span class="badge bg-success me-2">{{ tag.name }}</span>{{ tag.description }}</span>
				<span class="badge bg-secondary" v-if="tag.restricted"><i class="fa fa-lock"></i> restricted</span>
			</div>
		</router-link>
	</div>
	<empty-state icon="tags" title="No tags defined" v-else/>
</template>
