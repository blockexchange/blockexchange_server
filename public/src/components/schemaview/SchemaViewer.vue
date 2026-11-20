<script setup lang="ts">
import { onMounted, useTemplateRef } from 'vue';
import { createView } from "@buckaroobanzay/luanti-viewer";

import { create_schema_viewer_options } from "../../service/schema_viewer";
import type { Schema } from "../../api/types";

const props = defineProps<{
	schema: Schema;
}>();

const container = useTemplateRef<HTMLDivElement>("container");

// the viewer releases its resources once the element is removed from the document
onMounted(() => createView(container.value!, create_schema_viewer_options(props.schema)));
</script>

<template>
	<div ref="container" class="bx-viewer"></div>
</template>

<style scoped>
.bx-viewer {
	width: 100%;
	height: 70vh;
	min-height: 400px;
}
</style>
