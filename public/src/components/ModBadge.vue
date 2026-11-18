<script setup lang="ts">
defineProps<{
	name: string;
}>();
</script>

<template>
	<router-link :to="{path: '/search', query: {mod_name: name}}" class="badge text-bg-secondary text-decoration-none fw-normal">
		{{ name }}
	</router-link>
</template>
