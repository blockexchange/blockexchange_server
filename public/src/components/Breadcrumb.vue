<script setup lang="ts">
import type { BreadcrumbItem } from "../util/breadcrumbs";

defineProps<{
	items: BreadcrumbItem[];
}>();

function get_icon_class(item: BreadcrumbItem) {
	const cl: Record<string, boolean> = {};
	if (item.icon){
		cl.fa = true;
		cl["fa-" + item.icon] = true;
	}
	return cl;
}
</script>

<template>
	<nav>
		<ol class="breadcrumb text-bg-secondary" style="padding: 5px; border-radius: 5px;">
			<li class="breadcrumb-item" v-for="item in items" :key="item.name">
				<router-link class="link-light" :to="item.link" v-if="item.link">
					<i :class="get_icon_class(item)" v-if="item.icon"></i>
					{{ item.name }}
				</router-link>
				<span v-else>
					<i :class="get_icon_class(item)" v-if="item.icon"></i>
					{{ item.name }}
				</span>
			</li>
		</ol>
	</nav>
</template>
