<script setup lang="ts">
defineProps<{
	name: string;
}>();
</script>

<template>
	<router-link :to="{path: '/search', query: {tag_name: name}}">
		<span class="badge bg-success" style="margin-right: 5px;">
			<i class="fa fa-tag"></i>
			{{ name }}
		</span>
	</router-link>
</template>
