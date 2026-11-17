<script setup lang="ts">
defineProps<{
	name: string;
}>();
</script>

<template>
	<router-link :to="{path: '/search', query: {mod_name: name}}">
		<span class="badge bg-primary" style="margin-right: 5px;">
			<i class="fa fa-box-archive"></i>
			{{ name }}
		</span>
	</router-link>
</template>
