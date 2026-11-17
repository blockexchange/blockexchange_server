<script setup lang="ts">
import { get_schemamod_count } from "../../api/schemamod";
import BreadCrumb from "../Breadcrumb.vue";
import LoadingBlock from "../LoadingBlock.vue";
import { START, MODS } from "../../util/breadcrumbs";

const breadcrumb = [START, MODS];

const fetch_data = () => ({
	mods: get_schemamod_count()
});
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="container">
		<h5>Used mods</h5>
		<loading-block :fetch_data="fetch_data" v-slot="{ data }">
			<router-link
				:to="{path: '/search', query: {mod_name: mod.mod_name}}"
				class="btn btn-secondary m-1"
				v-for="mod in data.mods"
				:key="mod.mod_name">
				<i class="fa fa-box-archive"></i>
				{{ mod.mod_name }}
				<span class="badge bg-primary rounded-pill">{{ mod.count }}</span>
			</router-link>
		</loading-block>
	</div>
</template>
