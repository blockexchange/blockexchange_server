<script setup lang="ts">
import { computed } from 'vue';
import { useRoute, type RouteLocationRaw } from 'vue-router';

const props = defineProps<{
	current: number;
	pages: number;
}>();

const route = useRoute();

// windowed page list: 1 … 4 5 [6] 7 8 … 20
const items = computed(() => {
	const window = 2;
	const result: (number | "gap")[] = [];
	for (let i = 1; i <= props.pages; i++) {
		if (i == 1 || i == props.pages || Math.abs(i - props.current) <= window) {
			result.push(i);
		} else if (result[result.length - 1] != "gap") {
			result.push("gap");
		}
	}
	return result;
});

const get_route = (page: number): RouteLocationRaw => ({
	path: route.path,
	query: { ...route.query, page: page }
});
</script>

<template>
	<nav aria-label="Pages" v-if="pages > 1" class="mt-4">
		<ul class="pagination justify-content-center flex-wrap">
			<li class="page-item" :class="{disabled: current <= 1}">
				<router-link class="page-link" :to="get_route(current - 1)" aria-label="Previous page">
					<i class="fa fa-chevron-left"></i>
				</router-link>
			</li>
			<template v-for="(item, i) in items" :key="i">
				<li class="page-item disabled" v-if="item == 'gap'">
					<span class="page-link">…</span>
				</li>
				<li class="page-item" :class="{active: item == current}" v-else>
					<router-link class="page-link" :to="get_route(item)">{{ item }}</router-link>
				</li>
			</template>
			<li class="page-item" :class="{disabled: current >= pages}">
				<router-link class="page-link" :to="get_route(current + 1)" aria-label="Next page">
					<i class="fa fa-chevron-right"></i>
				</router-link>
			</li>
		</ul>
	</nav>
</template>
