<script setup lang="ts">
import { computed, ref, useTemplateRef } from 'vue';

import BreadCrumb from "../Breadcrumb.vue";
import LoadingBlock from "../LoadingBlock.vue";
import EmptyState from "../EmptyState.vue";
import SchemaViewer from "../schemaview/SchemaViewer.vue";
import { START, USERS, USER_DETAIL, SCHEMA_DETAIL, SCHEMA_VIEW_3D } from "../../util/breadcrumbs";

import { schema_search } from "../../api/schema";

const props = defineProps<{
	username: string;
	name: string;
}>();

const breadcrumb = computed(() => [
	START, USERS, USER_DETAIL(props.username),
	SCHEMA_DETAIL(props.username, props.name),
	SCHEMA_VIEW_3D(props.username, props.name)
]);

const detail_link = computed(() => `/schema/${props.username}/${props.name}`);

const viewer = useTemplateRef<InstanceType<typeof SchemaViewer>>("viewer");
const exporting = ref(false);
const export_error = ref<string | null>(null);

const export_gltf = async () => {
	if (!viewer.value || exporting.value) {
		return;
	}
	exporting.value = true;
	export_error.value = null;
	try {
		const blob = await viewer.value.export_gltf();
		const url = URL.createObjectURL(blob);
		const link = document.createElement("a");
		link.href = url;
		link.download = `${props.name}.glb`;
		link.click();
		setTimeout(() => URL.revokeObjectURL(url), 0);
	} catch (e) {
		console.error("glTF export failed", e);
		export_error.value = e instanceof Error ? e.message : String(e);
	} finally {
		exporting.value = false;
	}
};

const fetch_search_result = async () => {
	const r = await schema_search({
		schema_name: props.name,
		user_name: props.username
	});
	return r[0] ?? null;
};

const fetch_data = () => ({
	search_result: fetch_search_result()
});
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<loading-block :fetch_data="fetch_data" v-slot="{ data }">
		<empty-state icon="circle-question" title="Schematic not found" v-if="!data.search_result">
			<p><b>{{ name }}</b> by {{ username }} doesn't exist or was renamed.</p>
			<router-link :to="'/user/' + username">See schematics by {{ username }}</router-link>
		</empty-state>
		<template v-else>
			<div class="d-flex flex-wrap justify-content-between align-items-start gap-3 mb-4">
				<div class="flex-grow-1" style="min-width: 0;">
					<h1 class="h2 mb-1 text-break">{{ data.search_result.schema.name }}</h1>
					<div class="text-muted-bx">
						by <router-link :to="'/user/' + username">{{ username }}</router-link> · 3D view
					</div>
				</div>
				<div class="d-flex gap-2">
					<button type="button" class="btn btn-outline-primary" :disabled="exporting" @click="export_gltf"
						title="Download the schematic as binary glTF (.glb), e.g. for blender">
						<i class="fa" :class="exporting ? 'fa-spinner fa-spin' : 'fa-download'"></i>
						{{ exporting ? "Exporting..." : "Export glTF" }}
					</button>
					<router-link :to="detail_link" class="btn btn-outline-secondary">
						<i class="fa fa-arrow-left"></i> Back
					</router-link>
				</div>
			</div>
			<div class="alert alert-danger" v-if="export_error">
				glTF export failed: {{ export_error }}
			</div>
			<div class="card overflow-hidden">
				<schema-viewer ref="viewer" :schema="data.search_result.schema"/>
			</div>
			<p class="small text-muted-bx mt-2 mb-0">
				Click into the view to look around: <kbd>W</kbd> <kbd>A</kbd> <kbd>S</kbd> <kbd>D</kbd> move,
				<kbd>Space</kbd> / <kbd>Shift</kbd> up and down, <kbd>E</kbd> faster, <kbd>Esc</kbd> releases the mouse.
				Nodes without a known definition are shown as unknown blocks.
			</p>
		</template>
	</loading-block>
</template>
