<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import PageHeader from "../PageHeader.vue";
import { START, GIT_MEDIA_SOURCES } from "../../util/breadcrumbs";
import { has_permission } from "../../service/login";
import {
	get_git_media_sources,
	create_git_media_source,
	update_git_media_source,
	delete_git_media_source,
	run_git_media_source
} from "../../api/git_media_source";
import type { GitMediaSource, GitMediaSourceRequest } from "../../api/types";
import { get_cdb_vcs_packages, type CdbVcsPackage } from "../../api/cdb";
import format_time from "../../util/format_time";

const breadcrumb = [START, GIT_MEDIA_SOURCES];
const allowed = computed(() => has_permission("MEDIA"));

const list = ref<GitMediaSource[]>([]);
const error = ref("");
const busy = ref(false);
// ids of the sources with a manual run in progress
const running = ref<Set<string>>(new Set());

// id of the source being edited, null while adding a new one
const editing_id = ref<string | null>(null);
const empty_form = (): GitMediaSourceRequest => ({ git_url: "", cdb_package: "", media_license: "", license: "" });
const form = ref<GitMediaSourceRequest>(empty_form());

const form_valid = computed(() =>
	/^https?:\/\/.+/.test(form.value.git_url) &&
	form.value.git_url.length <= 256 &&
	(form.value.cdb_package == "" || /^[A-Za-z0-9_.-]+\/[A-Za-z0-9_.-]+$/.test(form.value.cdb_package)) &&
	form.value.media_license.length <= 64 &&
	form.value.license.length <= 64
);

// url autocompletion from the content db, loaded on first focus
const cdb_packages = ref<CdbVcsPackage[]>([]);
const max_suggestions = 50;

async function load_suggestions() {
	if (cdb_packages.value.length > 0) {
		return;
	}
	try {
		cdb_packages.value = await get_cdb_vcs_packages();
	} catch {
		// autocompletion is optional, the field stays a plain input
	}
}

const suggestions = computed(() => {
	const q = form.value.git_url.toLowerCase();
	const result: CdbVcsPackage[] = [];
	for (const p of cdb_packages.value) {
		if (!p.repo) {
			continue;
		}
		if (q == "" || p.repo.toLowerCase().includes(q) || p.title.toLowerCase().includes(q) || p.name.toLowerCase().includes(q) || p.author.toLowerCase().includes(q)) {
			result.push(p);
			if (result.length >= max_suggestions) {
				break;
			}
		}
	}
	return result;
});

// picking a suggestion fills in the contentdb package (used for the license lookup)
watch(() => form.value.git_url, url => {
	const p = cdb_packages.value.find(p => p.repo == url);
	if (p) {
		form.value.cdb_package = `${p.author}/${p.name}`;
	}
});

const error_message = (e: unknown) => (e as { message?: string })?.message || JSON.stringify(e);

async function update() {
	list.value = await get_git_media_sources();
}

function edit(s: GitMediaSource) {
	editing_id.value = s.id;
	form.value = { git_url: s.git_url, cdb_package: s.cdb_package, media_license: s.media_license, license: s.license };
	error.value = "";
}

function reset() {
	editing_id.value = null;
	form.value = empty_form();
}

async function save() {
	busy.value = true;
	error.value = "";
	try {
		if (editing_id.value) {
			await update_git_media_source(editing_id.value, form.value);
		} else {
			await create_git_media_source(form.value);
		}
		reset();
		await update();
	} catch (e) {
		error.value = `Saving failed: ${error_message(e)}`;
	} finally {
		busy.value = false;
	}
}

async function remove(s: GitMediaSource) {
	if (!confirm(`Delete the source "${s.git_url}" and all of its media files?`)) {
		return;
	}
	error.value = "";
	try {
		await delete_git_media_source(s.id);
		if (editing_id.value == s.id) {
			reset();
		}
		await update();
	} catch (e) {
		error.value = `Deleting "${s.git_url}" failed: ${error_message(e)}`;
	}
}

async function run(s: GitMediaSource) {
	running.value.add(s.id);
	error.value = "";
	try {
		const updated = await run_git_media_source(s.id);
		const i = list.value.findIndex(e => e.id == updated.id);
		if (i >= 0) {
			list.value[i] = updated;
		}
	} catch (e) {
		error.value = `Running "${s.git_url}" failed: ${error_message(e)}`;
	} finally {
		running.value.delete(s.id);
	}
}

onMounted(() => {
	if (allowed.value) {
		update().catch(e => error.value = `Loading failed: ${error_message(e)}`);
	}
});
</script>

<template>
	<page-header title="Git media sources"
		subtitle="Git repositories the textures and models are imported from, once a day or on demand with Run. Empty licenses are fetched from the ContentDB."
		:breadcrumb="breadcrumb"/>
	<div class="alert alert-warning" v-if="!allowed">
		<i class="fa fa-lock"></i> You need the <b>MEDIA</b> permission to manage git media sources.
	</div>
	<template v-else>
		<div class="alert alert-danger" v-if="error">{{ error }}</div>

		<div class="card mb-4">
			<div class="card-body">
				<h6 class="card-title">
					<template v-if="editing_id">
						<i class="fa fa-pen-to-square"></i> Edit source
					</template>
					<template v-else>
						<i class="fa fa-plus"></i> Add source
					</template>
				</h6>
				<form class="row g-2 align-items-end" @submit.prevent="save">
					<div class="col-lg-6">
						<label class="form-label small" for="gms-url">Git URL</label>
						<input id="gms-url" type="text" class="form-control" maxlength="256"
							placeholder="https://github.com/user/repo" v-model.trim="form.git_url"
							list="gms-url-list" autocomplete="off" @focus="load_suggestions">
						<datalist id="gms-url-list">
							<option v-for="p in suggestions" :key="p.author + '/' + p.name" :value="p.repo!">
								{{ p.title }} ({{ p.author }}/{{ p.name }}, {{ p.type }})
							</option>
						</datalist>
					</div>
					<div class="col-sm-4 col-lg-2">
						<label class="form-label small" for="gms-cdb">ContentDB package</label>
						<input id="gms-cdb" type="text" class="form-control" maxlength="128"
							placeholder="author/name" v-model.trim="form.cdb_package">
					</div>
					<div class="col-sm-4 col-lg-2">
						<label class="form-label small" for="gms-media-license">Media license</label>
						<input id="gms-media-license" type="text" class="form-control" maxlength="64"
							placeholder="from ContentDB" v-model.trim="form.media_license">
					</div>
					<div class="col-sm-4 col-lg-2">
						<label class="form-label small" for="gms-license">License</label>
						<input id="gms-license" type="text" class="form-control" maxlength="64"
							placeholder="from ContentDB" v-model.trim="form.license">
					</div>
					<div class="col-12 d-flex gap-2 justify-content-end">
						<button class="btn btn-outline-secondary" type="button" v-if="editing_id" @click="reset">
							Cancel
						</button>
						<button class="btn btn-success" type="submit" :disabled="!form_valid || busy">
							<span class="spinner-border spinner-border-sm" v-if="busy"></span>
							<i class="fa fa-floppy-disk" v-else></i>
							{{ editing_id ? "Save" : "Add" }}
						</button>
					</div>
				</form>
			</div>
		</div>

		<div class="card">
			<div class="table-responsive">
				<table class="table table-hover align-middle mb-0">
					<thead>
						<tr>
							<th>Repository</th>
							<th class="d-none d-md-table-cell">Licenses</th>
							<th>Last update</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						<tr v-if="list.length == 0">
							<td colspan="4" class="text-center text-muted-bx py-4">No sources defined</td>
						</tr>
						<tr v-for="s in list" :key="s.id" :class="{'table-active': editing_id == s.id}">
							<td class="font-monospace small text-break">
								<a :href="s.git_url" target="_blank" rel="noopener noreferrer">{{ s.git_url }}</a>
								<div class="text-muted-bx" v-if="s.cdb_package">
									<a :href="'https://content.luanti.org/packages/' + s.cdb_package + '/'" target="_blank" rel="noopener noreferrer" class="text-muted-bx">
										<i class="fa fa-box"></i> {{ s.cdb_package }}
									</a>
								</div>
							</td>
							<td class="d-none d-md-table-cell small">
								<div><span class="text-muted-bx">Media:</span> {{ s.media_license || "-" }}</div>
								<div><span class="text-muted-bx">Code:</span> {{ s.license || "-" }}</div>
							</td>
							<td class="small">
								<span class="text-muted-bx" v-if="!s.last_updated_timestamp">never</span>
								<template v-else>
									<span class="badge text-bg-success" v-if="s.last_updated_success"><i class="fa fa-check"></i> ok</span>
									<span class="badge text-bg-danger" v-else><i class="fa fa-xmark"></i> failed</span>
									<span class="text-muted-bx ms-1">{{ format_time(s.last_updated_timestamp) }}</span>
									<div class="text-muted-bx text-break" v-if="s.last_updated_message">{{ s.last_updated_message }}</div>
								</template>
							</td>
							<td class="text-end text-nowrap">
								<button type="button" class="btn btn-sm btn-outline-success me-1" @click="run(s)" :disabled="running.has(s.id)" :title="'Import media from ' + s.git_url + ' now'">
									<span class="spinner-border spinner-border-sm" v-if="running.has(s.id)"></span>
									<i class="fa fa-play" v-else></i>
									Run
								</button>
								<button type="button" class="btn btn-sm btn-outline-primary me-1" @click="edit(s)" :title="'Edit ' + s.git_url">
									<i class="fa fa-pen"></i>
									<span class="visually-hidden">Edit</span>
								</button>
								<button type="button" class="btn btn-sm btn-outline-danger" @click="remove(s)" :title="'Delete ' + s.git_url">
									<i class="fa fa-trash"></i>
									<span class="visually-hidden">Delete</span>
								</button>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	</template>
</template>
