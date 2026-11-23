<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';

import {
	get_schema_mod_dependencies,
	create_schema_mod_dependency,
	update_schema_mod_dependency,
	delete_schema_mod_dependency
} from "../../api/schema_mod_dependency";
import type { SchemaModDependency, SchemaModDependencyRequest } from "../../api/types";

// editor for the mod/game dependencies of a schematic, changes are saved right away
const props = defineProps<{
	schema_uid: string;
}>();

// emitted with the current list after every (re-)load
const emit = defineEmits<{
	change: [list: SchemaModDependency[]];
}>();

const list = ref<SchemaModDependency[]>([]);
const loading = ref(true);
const busy = ref(false);
const error = ref("");

// uid of the entry being edited, null while adding a new one
const editing_uid = ref<string | null>(null);
const empty_form = (): SchemaModDependencyRequest => ({ mod_type: "mod", git_url: "", version: "" });
const form = ref<SchemaModDependencyRequest>(empty_form());

const form_valid = computed(() =>
	form.value.git_url != "" &&
	form.value.git_url.length <= 256 &&
	form.value.version.length <= 64
);

const error_message = (e: unknown) => (e as { message?: string })?.message || JSON.stringify(e);

async function update() {
	list.value = await get_schema_mod_dependencies(props.schema_uid);
	emit("change", list.value);
}

async function load() {
	loading.value = true;
	error.value = "";
	try {
		await update();
	} catch (e) {
		error.value = `Loading failed: ${error_message(e)}`;
	} finally {
		loading.value = false;
	}
}

function edit(d: SchemaModDependency) {
	editing_uid.value = d.uid;
	form.value = { mod_type: d.mod_type, git_url: d.git_url, version: d.version };
	error.value = "";
}

function reset() {
	editing_uid.value = null;
	form.value = empty_form();
}

async function save() {
	busy.value = true;
	error.value = "";
	try {
		if (editing_uid.value) {
			await update_schema_mod_dependency(props.schema_uid, editing_uid.value, form.value);
		} else {
			await create_schema_mod_dependency(props.schema_uid, form.value);
		}
		reset();
		await update();
	} catch (e) {
		error.value = `Saving failed: ${error_message(e)}`;
	} finally {
		busy.value = false;
	}
}

async function remove(d: SchemaModDependency) {
	if (!confirm(`Remove the dependency "${d.git_url}"?`)) {
		return;
	}
	error.value = "";
	try {
		await delete_schema_mod_dependency(props.schema_uid, d.uid);
		if (editing_uid.value == d.uid) {
			reset();
		}
		await update();
	} catch (e) {
		error.value = `Removing "${d.git_url}" failed: ${error_message(e)}`;
	}
}

onMounted(load);
watch(() => props.schema_uid, () => {
	reset();
	load();
});
</script>

<template>
	<div>
		<p class="small text-muted-bx">
			Git repositories of the mods or game this schematic depends on.
			Changes here are saved immediately.
		</p>

		<div class="alert alert-danger py-2 small" v-if="error">{{ error }}</div>

		<div class="table-responsive mb-3">
			<table class="table table-sm align-middle mb-0">
				<thead>
					<tr>
						<th>Type</th>
						<th>Git URL</th>
						<th>Version</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					<tr v-if="loading">
						<td colspan="4" class="text-center py-3">
							<span class="spinner-border spinner-border-sm"></span>
						</td>
					</tr>
					<tr v-else-if="list.length == 0">
						<td colspan="4" class="text-center text-muted-bx py-3 small">No dependencies defined</td>
					</tr>
					<tr v-for="d in list" :key="d.uid" :class="{'table-active': editing_uid == d.uid}">
						<td>
							<span class="badge" :class="d.mod_type == 'game' ? 'text-bg-primary' : 'text-bg-secondary'">
								<i class="fa" :class="d.mod_type == 'game' ? 'fa-gamepad' : 'fa-puzzle-piece'"></i>
								{{ d.mod_type }}
							</span>
						</td>
						<td class="font-monospace small text-break">
							<a :href="d.git_url" target="_blank" rel="noopener noreferrer" v-if="/^https?:\/\//.test(d.git_url)">{{ d.git_url }}</a>
							<template v-else>{{ d.git_url }}</template>
						</td>
						<td class="font-monospace small">{{ d.version }}</td>
						<td class="text-end text-nowrap">
							<button type="button" class="btn btn-sm btn-outline-primary me-1" @click="edit(d)" :title="'Edit ' + d.git_url">
								<i class="fa fa-pen"></i>
								<span class="visually-hidden">Edit</span>
							</button>
							<button type="button" class="btn btn-sm btn-outline-danger" @click="remove(d)" :title="'Remove ' + d.git_url">
								<i class="fa fa-trash"></i>
								<span class="visually-hidden">Remove</span>
							</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>

		<h6 class="small fw-semibold">
			<template v-if="editing_uid"><i class="fa fa-pen-to-square"></i> Edit dependency</template>
			<template v-else><i class="fa fa-plus"></i> Add dependency</template>
		</h6>
		<form class="row g-2 align-items-end" @submit.prevent="save">
			<div class="col-sm-3">
				<label class="form-label small" for="dep-type">Type</label>
				<select id="dep-type" class="form-select form-select-sm" v-model="form.mod_type">
					<option value="mod">mod</option>
					<option value="game">game</option>
				</select>
			</div>
			<div class="col-sm-6">
				<label class="form-label small" for="dep-url">Git URL</label>
				<input id="dep-url" type="text" class="form-control form-control-sm" maxlength="256"
					placeholder="https://github.com/user/repo.git" v-model.trim="form.git_url">
			</div>
			<div class="col-sm-3">
				<label class="form-label small" for="dep-version">Version</label>
				<input id="dep-version" type="text" class="form-control form-control-sm" maxlength="64"
					placeholder="latest" v-model.trim="form.version">
			</div>
			<div class="col-12 d-flex gap-2 justify-content-end">
				<button class="btn btn-sm btn-outline-secondary" type="button" v-if="editing_uid" @click="reset">
					Cancel
				</button>
				<button class="btn btn-sm btn-success" type="submit" :disabled="!form_valid || busy">
					<span class="spinner-border spinner-border-sm" v-if="busy"></span>
					<i class="fa fa-floppy-disk" v-else></i>
					{{ editing_uid ? "Save" : "Add" }}
				</button>
			</div>
		</form>
	</div>
</template>
