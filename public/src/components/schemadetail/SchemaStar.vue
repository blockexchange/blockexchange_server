<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { is_logged_in } from "../../service/login";
import { get_schema_star, star_schema, unstar_schema, count_schema_stars } from "../../api/schema_star";

const props = defineProps<{
	schema_uid: string;
}>();

// star count of the schema
const stars = defineModel<number>("stars", { required: true });

const schema_star = ref<unknown>(null);
const logged_in = computed(is_logged_in);

async function update_star() {
	if (is_logged_in()) {
		schema_star.value = await get_schema_star(props.schema_uid);
		stars.value = await count_schema_stars(props.schema_uid);
	}
}

async function star() {
	await star_schema(props.schema_uid);
	await update_star();
}

async function unstar() {
	await unstar_schema(props.schema_uid);
	await update_star();
}

onMounted(update_star);
</script>

<template>
	<button type="button" class="btn"
		:class="schema_star ? 'btn-warning' : 'btn-outline-secondary'"
		:disabled="!logged_in"
		:title="logged_in ? (schema_star ? 'Remove your star' : 'Star this schematic') : 'Log in to star schematics'"
		@click="schema_star ? unstar() : star()">
		<i :class="schema_star ? 'fa-solid fa-star' : 'fa-regular fa-star bx-star'"></i>
		{{ schema_star ? 'Starred' : 'Star' }}
		<span class="badge rounded-pill ms-1" :class="schema_star ? 'text-bg-dark' : 'text-bg-secondary'">{{ stars }}</span>
	</button>
</template>
