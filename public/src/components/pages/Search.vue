<script setup lang="ts">
import BreadCrumb from "../Breadcrumb.vue";
import SchemaSearch from "../SchemaSearch.vue";
import { START, SEARCH } from "../../util/breadcrumbs";

const breadcrumb = [START, SEARCH];
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<schema-search/>
</template>
