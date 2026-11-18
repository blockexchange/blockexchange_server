<script setup lang="ts" generic="T">
import { computed, onMounted, ref, watch, type Ref } from 'vue';
import { useRoute } from 'vue-router';
import Pagination from './Pagination.vue';
import LoadingSpinner from './LoadingSpinner.vue';

const props = defineProps<{
	fetch_entries: (limit: number, offset: number) => Promise<T[]>;
	count_entries: () => Promise<number>;
	per_page: number;
}>();

defineSlots<{
	body(props: { list: T[]; total: number }): unknown;
}>();

const route = useRoute();

const list = ref([]) as Ref<T[]>;
const busy = ref(false);
const total_count = ref(0);
const current_page = ref(1);

const pages = computed(() => Math.ceil(total_count.value / props.per_page));

async function update() {
	busy.value = true;
	list.value = [];
	current_page.value = +(route.query.page || 1);
	try {
		total_count.value = await props.count_entries();
		list.value = await props.fetch_entries(props.per_page, (current_page.value-1)*props.per_page);
	} finally {
		busy.value = false;
	}
}

onMounted(update);
watch(() => route.fullPath, update);
</script>

<template>
	<loading-spinner v-if="busy"/>
	<template v-else>
		<slot name="body" :list="list" :total="total_count"></slot>
		<pagination :current="current_page" :pages="pages"/>
	</template>
</template>
