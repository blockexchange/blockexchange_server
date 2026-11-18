<script setup lang="ts">
import NavBar from './components/NavBar.vue';
import ErrorToast from './components/ErrorToast.vue';
import AppFooter from './components/AppFooter.vue';
</script>

<template>
	<nav-bar/>
	<error-toast/>
	<main class="container-xl py-4">
		<router-view/>
	</main>
	<app-footer/>
</template>
