<script setup lang="ts">
import { computed, reactive, ref, useTemplateRef } from 'vue';

import PageHeader from "../PageHeader.vue";
import EmptyState from "../EmptyState.vue";
import { START, SCHEMA_IMPORT } from "../../util/breadcrumbs";
import { import_schematic } from "../../api/import";
import type { ErrorResponse } from "../../api/types";
import { get_username } from "../../service/login";
import store, { type ImportResult } from "../../store/import";

const breadcrumb = [START, SCHEMA_IMPORT];
const username = computed(get_username);
const upload = useTemplateRef<HTMLInputElement>("upload");
const dragging = ref(false);

async function import_file(file: File) {
	const res = reactive<ImportResult>({
		filename: file.name,
		busy: true,
		schema: null,
		error: null
	});
	store.results.unshift(res);

	try {
		const buf = await file.arrayBuffer();
		res.schema = await import_schematic(buf, file.name);
	} catch (e) {
		res.error = e as ErrorResponse;
	} finally {
		res.busy = false;
	}
}

function import_files(files: FileList | null | undefined) {
	for (const file of files ?? []) {
		import_file(file);
	}
}

function on_upload() {
	import_files(upload.value?.files);
	if (upload.value) {
		// allow selecting the same file again
		upload.value.value = "";
	}
}

function on_drop(e: DragEvent) {
	dragging.value = false;
	import_files(e.dataTransfer?.files);
}
</script>

<template>
	<page-header
		title="Import schematic files"
		subtitle="Upload WorldEdit (.we) or Block exchange (.zip) files from your computer. To upload directly from a world, use /bx_save in-game."
		:breadcrumb="breadcrumb"/>

	<empty-state icon="lock" title="Log in to import schematics" v-if="!username">
		<router-link to="/login" class="btn btn-success mt-2">Log in</router-link>
	</empty-state>

	<div class="row justify-content-center" v-else>
		<div class="col-lg-8">
			<label
				class="bx-dropzone d-block mb-4"
				:class="{active: dragging}"
				@dragover.prevent="dragging = true"
				@dragleave.prevent="dragging = false"
				@drop.prevent="on_drop">
				<i class="fa fa-file-arrow-up fa-2x mb-2 text-muted-bx"></i>
				<div class="fw-semibold">Drop files here or click to choose</div>
				<div class="small text-muted-bx">.we and .zip, several files at once are fine</div>
				<input ref="upload" type="file" class="visually-hidden" @change="on_upload" multiple accept=".we,.zip"/>
			</label>

			<div class="list-group" v-if="store.results.length">
				<div class="list-group-item d-flex align-items-center gap-3" v-for="(res, i) in store.results" :key="i">
					<span class="spinner-border spinner-border-sm text-muted-bx" v-if="res.busy"></span>
					<i class="fa fa-circle-check text-success" v-else-if="res.schema"></i>
					<i class="fa fa-circle-xmark text-danger" v-else></i>
					<div class="flex-grow-1 text-break">
						<div>{{ res.filename }}</div>
						<div class="small text-danger" v-if="res.error">
							{{ res.error.message || 'Import failed' }}
						</div>
						<div class="small text-muted-bx" v-else-if="res.busy">Importing…</div>
					</div>
					<router-link class="btn btn-sm btn-outline-success" :to="'/schema/' + username + '/' + res.schema.name" v-if="res.schema">
						Open {{ res.schema.name }}
					</router-link>
				</div>
			</div>
		</div>
	</div>
</template>
