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
	<div class="toast-container top-50 start-50 translate-middle">
		<div class="toast show" v-if="store.message">
			<div class="toast-header text-bg-danger">
				<strong class="me-auto">
					<i class="fa-solid fa-triangle-exclamation"></i>
					Error
				</strong>
				<small>
					{{ store.title }}
					<span v-if="store.status">Code: {{ store.status }}</span>
				</small>
				<button type="button" class="btn-close" @click="close"></button>
			</div>
			<div class="toast-body">
				<p>
					<b>Message:</b> {{ store.message }}
				</p>
				<p>
					<b>URL:</b> {{ store.url }}
				</p>
			</div>
		</div>
	</div>
</template>
