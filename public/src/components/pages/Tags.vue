<script setup lang="ts">
import BreadCrumb from "../Breadcrumb.vue";
import { START, TAGS } from "../../util/breadcrumbs";

const breadcrumb = [START, TAGS];
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
</template>
