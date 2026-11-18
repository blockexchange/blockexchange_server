<script setup lang="ts">
import type { BreadcrumbItem } from "../util/breadcrumbs";

defineProps<{
	items: BreadcrumbItem[];
}>();
</script>

<template>
	<nav aria-label="breadcrumb" v-if="items.length > 1">
		<ol class="breadcrumb">
			<li class="breadcrumb-item"
				v-for="(item, i) in items"
				:key="item.name"
				:class="{active: i == items.length - 1}"
				:aria-current="i == items.length - 1 ? 'page' : undefined">
				<router-link :to="item.link" v-if="item.link && i < items.length - 1">
					<i :class="['fa', 'fa-' + item.icon]" v-if="item.icon && i == 0"></i>
					{{ item.name }}
				</router-link>
				<span v-else>{{ item.name }}</span>
			</li>
		</ol>
	</nav>
</template>
