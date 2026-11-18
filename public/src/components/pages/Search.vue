<script setup lang="ts">
import PageHeader from "../PageHeader.vue";
import SchemaSearch from "../SchemaSearch.vue";
import { START, SEARCH } from "../../util/breadcrumbs";

const breadcrumb = [START, SEARCH];
</script>

<template>
	<page-header
		title="Browse schematics"
		subtitle="Find builds shared by the community. Every schematic can be placed in-game with the blockexchange mod or downloaded as a file."
		:breadcrumb="breadcrumb"/>
	<schema-search/>
</template>
