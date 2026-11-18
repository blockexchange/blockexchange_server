<script setup lang="ts">
import { ref } from 'vue';
import { get_schemamod_count } from "../../api/schemamod";
import PageHeader from "../PageHeader.vue";
import LoadingBlock from "../LoadingBlock.vue";
import EmptyState from "../EmptyState.vue";
import { START, MODS } from "../../util/breadcrumbs";

const breadcrumb = [START, MODS];
const filter = ref("");

const fetch_data = () => ({
	mods: get_schemamod_count()
});
</script>

<template>
	<page-header
		title="Mods"
		subtitle="The mods used by shared schematics. A schematic only places correctly if your world has the same mods installed. Pick a mod to see the schematics that use it."
		:breadcrumb="breadcrumb"/>
	<loading-block :fetch_data="fetch_data" v-slot="{ data }">
		<div class="input-group mb-3" style="max-width: 24rem;">
			<span class="input-group-text"><i class="fa fa-filter"></i></span>
			<input type="search" class="form-control" placeholder="Filter mods" aria-label="Filter mods" v-model="filter"/>
		</div>
		<div class="d-flex flex-wrap gap-2">
			<router-link
				:to="{path: '/search', query: {mod_name: mod.mod_name}}"
				class="btn btn-outline-secondary btn-sm"
				v-for="mod in data.mods.filter(m => m.mod_name.toLowerCase().includes(filter.toLowerCase()))"
				:key="mod.mod_name">
				{{ mod.mod_name }}
				<span class="badge text-bg-secondary rounded-pill ms-1">{{ mod.count }}</span>
			</router-link>
		</div>
		<empty-state icon="box-archive" title="No mods match your filter"
			v-if="!data.mods.some(m => m.mod_name.toLowerCase().includes(filter.toLowerCase()))"/>
	</loading-block>
</template>
