<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
	text: string;
}>();

const copied = ref(false);
let timer: ReturnType<typeof setTimeout> | undefined;

async function copy() {
	await navigator.clipboard.writeText(props.text);
	copied.value = true;
	clearTimeout(timer);
	timer = setTimeout(() => copied.value = false, 1500);
}
</script>

<template>
	<span class="bx-cmd">
		<code>{{ text }}</code>
		<button type="button" :class="{copied}" @click="copy" :title="copied ? 'Copied' : 'Copy to clipboard'">
			<i class="fa" :class="copied ? 'fa-check' : 'fa-copy'"></i>
			<span class="visually-hidden">Copy</span>
		</button>
	</span>
</template>
