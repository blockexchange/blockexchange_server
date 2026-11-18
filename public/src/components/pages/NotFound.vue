<script setup lang="ts">
import EmptyState from "../EmptyState.vue";
</script>

<template>
	<empty-state icon="map-location-dot" title="Page not found">
		<p>This page doesn't exist, or the link is outdated.</p>
		<div class="d-flex justify-content-center gap-2">
			<router-link to="/" class="btn btn-outline-light">Home</router-link>
			<router-link to="/search" class="btn btn-success">Browse schematics</router-link>
		</div>
	</empty-state>
</template>
