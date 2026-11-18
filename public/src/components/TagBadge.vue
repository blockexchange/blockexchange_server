<script setup lang="ts">
defineProps<{
	name: string;
}>();
</script>

<template>
	<router-link :to="{path: '/search', query: {tag_name: name}}" class="badge bg-success text-decoration-none fw-normal">
		{{ name }}
	</router-link>
</template>
