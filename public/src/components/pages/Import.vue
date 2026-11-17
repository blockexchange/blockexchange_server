<script setup lang="ts">
import { computed, reactive, useTemplateRef } from 'vue';

import BreadCrumb from "../Breadcrumb.vue";
import { START, SCHEMA_IMPORT } from "../../util/breadcrumbs";
import { import_schematic } from "../../api/import";
import type { ErrorResponse } from "../../api/types";
import { get_username } from "../../service/login";
import store, { type ImportResult } from "../../store/import";

const breadcrumb = [START, SCHEMA_IMPORT];
const username = computed(get_username);
const upload = useTemplateRef<HTMLInputElement>("upload");

function import_file(file: File) {
	const res = reactive<ImportResult>({
		filename: file.name,
		busy: true,
		schema: null,
		error: null
	});
	store.results.push(res);

	file.arrayBuffer()
	.then(buf => import_schematic(buf, file.name))
	.then(s => res.schema = s)
	.catch((e: ErrorResponse) => res.error = e)
	.finally(() => res.busy = false);
}

function on_upload() {
	store.results = [];
	for (const file of upload.value?.files ?? []) {
		import_file(file);
	}
}
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="row">
		<div class="col-md-4"></div>
		<div class="col-md-4">
			<h5>
				Upload Worldedit or Blockexchange schematics
			</h5>
			<input ref="upload" type="file" class="form-control" @change="on_upload" multiple accept=".we,.zip"/>
			<hr>
			<ul>
				<li v-for="res in store.results" :key="res.filename">
					{{ res.filename }}
					<i class="fa fa-spinner fa-spin" v-if="res.busy"></i>
					<i class="fa fa-check" v-if="res.schema"></i>
					<i class="fa fa-times" v-if="res.error"></i>
					<span class="badge bg-danger" v-if="res.error && res.error.message">
						{{ res.error.message }}
					</span>
					<router-link :to="'/schema/' + username + '/' + res.schema.name" v-if="res.schema">
						{{ res.schema.name }}
					</router-link>
				</li>
			</ul>
		</div>
		<div class="col-md-4"></div>
	</div>
</template>
