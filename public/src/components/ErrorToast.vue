<script setup lang="ts">
import store from "../store/error_toast";

function close() {
	store.message = "";
	store.title = "";
	store.url = "";
	store.status = 0;
}
</script>

<template>
	<div class="toast-container position-fixed bottom-0 end-0 p-3" style="z-index: 1100;">
		<div class="toast show border-danger" role="alert" aria-live="assertive" v-if="store.message">
			<div class="toast-header text-bg-danger">
				<i class="fa-solid fa-triangle-exclamation me-2"></i>
				<strong class="me-auto">Something went wrong</strong>
				<small v-if="store.status">HTTP {{ store.status }}</small>
				<button type="button" class="btn-close btn-close-white ms-2" aria-label="Close" @click="close"></button>
			</div>
			<div class="toast-body">
				<p class="mb-1">{{ store.message }}</p>
				<small class="text-muted-bx text-break" v-if="store.url">{{ store.title }} {{ store.url }}</small>
			</div>
		</div>
	</div>
</template>
