<script setup lang="ts">
import { onMounted, useTemplateRef } from 'vue';
import { createView, type Viewer } from "@buckaroobanzay/luanti-viewer";

import { create_schema_viewer_options, get_schema_bounds } from "../../service/schema_viewer";
import type { Schema } from "../../api/types";

const props = defineProps<{
	schema: Schema;
}>();

const container = useTemplateRef<HTMLDivElement>("container");
let viewer: Viewer | null = null;

// the viewer releases its resources once the element is removed from the document
onMounted(() => {
	viewer = createView(container.value!, create_schema_viewer_options(props.schema));
});

// exports the whole schematic (not only the loaded mapblocks) as binary glTF
const export_gltf = async (): Promise<Blob> => {
	if (!viewer) {
		throw new Error("viewer not ready");
	}
	const bounds = get_schema_bounds(props.schema);
	return viewer.exportGLTF({ minPos: bounds.min, maxPos: bounds.max });
};

defineExpose({ export_gltf });
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
