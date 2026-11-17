<script setup lang="ts">
import { ref } from 'vue';

const props = defineProps<{
	text: string;
}>();

const copied = ref(false);

function copy() {
	copied.value = true;
	navigator.clipboard.writeText(props.text);
}
</script>

<template>
	<span class="badge bg-secondary" @click="copy">
		{{ text }}
		&nbsp;
		<i class="fa fa-copy" :style="{'color': copied ? 'green' : undefined}"></i>
	</span>
</template>
