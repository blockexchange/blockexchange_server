<script setup lang="ts">
import BreadCrumb from "./Breadcrumb.vue";
import type { BreadcrumbItem } from "../util/breadcrumbs";

defineProps<{
	title: string;
	subtitle?: string;
	breadcrumb?: BreadcrumbItem[];
}>();

defineSlots<{
	actions?(): unknown;
	default?(): unknown;
}>();
</script>

<template>
	<header class="bx-page-header">
		<bread-crumb :items="breadcrumb" v-if="breadcrumb"/>
		<div class="d-flex flex-wrap gap-3 justify-content-between align-items-end">
			<div>
				<h1>{{ title }}</h1>
				<p class="lead" v-if="subtitle">{{ subtitle }}</p>
				<slot></slot>
			</div>
			<div class="d-flex flex-wrap gap-2" v-if="$slots.actions">
				<slot name="actions"></slot>
			</div>
		</div>
	</header>
</template>
