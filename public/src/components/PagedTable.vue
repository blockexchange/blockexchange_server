<script setup lang="ts" generic="T">
import { computed, onMounted, ref, watch, type Ref } from 'vue';
import { useRoute, type RouteLocationRaw } from 'vue-router';

// the table classes are passed as attributes
defineOptions({ inheritAttrs: false });

const props = defineProps<{
	fetch_entries: (limit: number, offset: number) => Promise<T[]>;
	count_entries: () => Promise<number>;
}>();

defineSlots<{
	header(): unknown;
	body(props: { list: T[] }): unknown;
}>();

const route = useRoute();

const per_page = 20;
const list = ref([]) as Ref<T[]>;
const busy = ref(false);
const total_count = ref(0);
const current_page = ref(1);

const pages = computed(() => Math.ceil(total_count.value / per_page));

async function update() {
	busy.value = true;
	current_page.value = +(route.query.page || 1);
	try {
		total_count.value = await props.count_entries();
		list.value = await props.fetch_entries(per_page, (current_page.value-1)*per_page);
	} finally {
		busy.value = false;
	}
}

function get_route(page: number): RouteLocationRaw {
	return {
		path: route.path,
		query: { ...route.query, page: page }
	};
}

onMounted(update);
watch(() => route.fullPath, update);
</script>

<template>
	<table v-bind="$attrs">
		<thead>
			<slot name="header"></slot>
		</thead>
		<tbody>
			<slot name="body" :list="list"></slot>
		</tbody>
	</table>
	<div class="btn-group">
		<router-link
			:to="get_route(i)"
			v-for="i in pages"
			:key="i"
			class="btn btn-xs btn-secondary"
			:class="{disabled: current_page == i}">
			{{ i }}
		</router-link>
	</div>
</template>
