<script setup lang="ts" generic="T">
import { computed, onMounted, ref, watch, type Ref } from 'vue';
import { useRoute, type RouteLocationRaw } from 'vue-router';

const props = defineProps<{
	fetch_entries: (limit: number, offset: number) => Promise<T[]>;
	count_entries: () => Promise<number>;
	per_page: number;
}>();

defineSlots<{
	body(props: { list: T[] }): unknown;
	busy?(): unknown;
}>();

const route = useRoute();

const list = ref([]) as Ref<T[]>;
const busy = ref(false);
const total_count = ref(0);
const current_page = ref(1);

const pages = computed(() => Math.ceil(total_count.value / props.per_page));

function update() {
	busy.value = true;
	list.value = [];
	current_page.value = +(route.query.page || 1);
	props.count_entries()
	.then(c => total_count.value = c)
	.then(() => props.fetch_entries(props.per_page, (current_page.value-1)*props.per_page))
	.then(l => list.value = l)
	.then(() => busy.value = false);
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
	<slot name="body" :list="list"></slot>
	<slot name="busy" v-if="busy">
		<div class="alert alert-primary">
			<i class="fa fa-spinner fa-spin"></i> Loading...
		</div>
	</slot>
	<div class="btn-group" v-if="!busy">
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
