<script setup lang="ts" generic="T extends Record<string, unknown>">
import { onMounted, ref, type Ref } from 'vue';
import LoadingSpinner from "./LoadingSpinner.vue";

// fetches all promises returned by "fetch_data" and passes the results as "data" to the default slot
const props = defineProps<{
	fetch_data: () => { [K in keyof T]: Promise<T[K]> };
}>();

defineSlots<{
	default(props: { data: T; update_data: () => void }): unknown;
}>();

const busy = ref(false);
const data = ref(null) as Ref<T | null>;

function update_data() {
	busy.value = true;
	const promise_map = props.fetch_data();
	const keys = Object.keys(promise_map) as (keyof T)[];
	Promise.all(keys.map(k => promise_map[k]))
	.then(values => {
		const result = {} as T;
		keys.forEach((k, i) => result[k] = values[i] as T[typeof k]);
		data.value = result;
		busy.value = false;
	});
}

onMounted(update_data);
</script>

<template>
	<loading-spinner v-if="busy"/>
	<slot :data="data" :update_data="update_data" v-if="data"></slot>
</template>
