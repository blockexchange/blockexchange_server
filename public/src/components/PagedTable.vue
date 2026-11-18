<script setup lang="ts" generic="T">
import { computed, onMounted, ref, watch, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import Pagination from './Pagination.vue';
import LoadingSpinner from './LoadingSpinner.vue';

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

onMounted(update);
watch(() => route.fullPath, update);
</script>

<template>
	<div class="table-responsive">
		<table v-bind="$attrs">
			<thead>
				<slot name="header"></slot>
			</thead>
			<tbody>
				<slot name="body" :list="list"></slot>
			</tbody>
		</table>
	</div>
	<loading-spinner v-if="busy && list.length == 0"/>
	<pagination :current="current_page" :pages="pages" v-if="!busy"/>
</template>
