<script setup lang="ts">
import NavBar from './components/NavBar.vue';
import ErrorToast from './components/ErrorToast.vue';
</script>

<template>
	<div>
		<nav-bar/>
		<error-toast/>
		<div class="container-fluid">
			<br>
			<router-view/>
		</div>
	</div>
</template>
