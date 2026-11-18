<script setup lang="ts">
withDefaults(defineProps<{
	icon?: string;
	title: string;
}>(), {
	icon: "box-open"
});
</script>

<template>
	<div class="bx-empty">
		<i :class="['fa', 'fa-' + icon]"></i>
		<h2 class="h5 text-light">{{ title }}</h2>
		<slot></slot>
	</div>
</template>
