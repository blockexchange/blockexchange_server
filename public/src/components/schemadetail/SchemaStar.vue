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

function update_star() {
	if (is_logged_in()) {
		get_schema_star(props.schema_uid)
		.then(s => schema_star.value = s)
		.then(() => count_schema_stars(props.schema_uid))
		.then(count => stars.value = count);
	}
}

const star = () => star_schema(props.schema_uid).then(update_star);
const unstar = () => unstar_schema(props.schema_uid).then(update_star);

onMounted(update_star);
</script>

<template>
	<button class="btn btn-outline-primary" :disabled="!logged_in" @click="schema_star ? unstar() : star()">
		<i class="fa fa-star" :style="{color: schema_star ? 'yellow' : ''}"></i>
		<span class="badge bg-secondary rouded-pill">{{ stars }}</span>
		{{ schema_star ? 'Unstar' : 'Star' }}
	</button>
</template>
