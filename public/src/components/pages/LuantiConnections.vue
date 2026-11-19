<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import PageHeader from "../PageHeader.vue";
import { START, LUANTI_CONNECTIONS } from "../../util/breadcrumbs";
import { has_permission } from "../../service/login";
import {
	get_luanti_media_connections,
	create_luanti_media_connection,
	update_luanti_media_connection,
	delete_luanti_media_connection,
	run_luanti_media_connection
} from "../../api/luanti_media_connection";
import type { LuantiMediaConnection, LuantiMediaConnectionRequest } from "../../api/types";
import format_time from "../../util/format_time";

const breadcrumb = [START, LUANTI_CONNECTIONS];
const allowed = computed(() => has_permission("MEDIA"));

const list = ref<LuantiMediaConnection[]>([]);
const error = ref("");
const busy = ref(false);
// ids of the connections with a manual run in progress
const running = ref<Set<string>>(new Set());

// id of the connection being edited, null while adding a new one
const editing_id = ref<string | null>(null);
const empty_form = (): LuantiMediaConnectionRequest => ({ hostname: "", port: 30000, username: "", password: "" });
const form = ref<LuantiMediaConnectionRequest>(empty_form());

const form_valid = computed(() =>
	form.value.hostname != "" &&
	form.value.username != "" &&
	form.value.port >= 1 && form.value.port <= 65535 &&
	(editing_id.value != null || form.value.password != "")
);

const error_message = (e: unknown) => (e as { message?: string })?.message || JSON.stringify(e);

async function update() {
	list.value = await get_luanti_media_connections();
}

function edit(c: LuantiMediaConnection) {
	editing_id.value = c.id;
	form.value = { hostname: c.hostname, port: c.port, username: c.username, password: "" };
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
			await update_luanti_media_connection(editing_id.value, form.value);
		} else {
			await create_luanti_media_connection(form.value);
		}
		reset();
		await update();
	} catch (e) {
		error.value = `Saving failed: ${error_message(e)}`;
	} finally {
		busy.value = false;
	}
}

async function remove(c: LuantiMediaConnection) {
	if (!confirm(`Delete the connection to "${c.hostname}:${c.port}"?`)) {
		return;
	}
	error.value = "";
	try {
		await delete_luanti_media_connection(c.id);
		if (editing_id.value == c.id) {
			reset();
		}
		await update();
	} catch (e) {
		error.value = `Deleting "${c.hostname}:${c.port}" failed: ${error_message(e)}`;
	}
}

async function run(c: LuantiMediaConnection) {
	running.value.add(c.id);
	error.value = "";
	try {
		const updated = await run_luanti_media_connection(c.id);
		const i = list.value.findIndex(e => e.id == updated.id);
		if (i >= 0) {
			list.value[i] = updated;
		}
	} catch (e) {
		error.value = `Running "${c.hostname}:${c.port}" failed: ${error_message(e)}`;
	} finally {
		running.value.delete(c.id);
	}
}

onMounted(() => {
	if (allowed.value) {
		update().catch(e => error.value = `Loading failed: ${error_message(e)}`);
	}
});
</script>

<template>
	<page-header title="Luanti connections" subtitle="Luanti servers that node definitions and media files are fetched from, once a day or on demand with Run." :breadcrumb="breadcrumb"/>
	<div class="alert alert-warning" v-if="!allowed">
		<i class="fa fa-lock"></i> You need the <b>MEDIA</b> permission to manage Luanti connections.
	</div>
	<template v-else>
		<div class="alert alert-danger" v-if="error">{{ error }}</div>

		<div class="card mb-4">
			<div class="card-body">
				<h6 class="card-title">
					<template v-if="editing_id">
						<i class="fa fa-pen-to-square"></i> Edit connection
					</template>
					<template v-else>
						<i class="fa fa-plus"></i> Add connection
					</template>
				</h6>
				<form class="row g-2 align-items-end" @submit.prevent="save">
					<div class="col-sm-8 col-lg-4">
						<label class="form-label small" for="lc-hostname">Hostname</label>
						<input id="lc-hostname" class="form-control" placeholder="e.g. luanti.example.com" v-model.trim="form.hostname">
					</div>
					<div class="col-sm-4 col-lg-2">
						<label class="form-label small" for="lc-port">Port</label>
						<input id="lc-port" class="form-control" type="number" min="1" max="65535" v-model.number="form.port">
					</div>
					<div class="col-sm-6 col-lg-3">
						<label class="form-label small" for="lc-username">Username</label>
						<input id="lc-username" class="form-control" autocomplete="off" v-model.trim="form.username">
					</div>
					<div class="col-sm-6 col-lg-3">
						<label class="form-label small" for="lc-password">Password</label>
						<input id="lc-password" class="form-control" type="password" autocomplete="new-password"
							:placeholder="editing_id ? 'unchanged' : ''" v-model="form.password">
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
							<th>Server</th>
							<th class="d-none d-sm-table-cell">Username</th>
							<th>Last update</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						<tr v-if="list.length == 0">
							<td colspan="4" class="text-center text-muted-bx py-4">No connections defined</td>
						</tr>
						<tr v-for="c in list" :key="c.id" :class="{'table-active': editing_id == c.id}">
							<td class="font-monospace">{{ c.hostname }}:{{ c.port }}</td>
							<td class="d-none d-sm-table-cell">{{ c.username }}</td>
							<td class="small">
								<span class="text-muted-bx" v-if="!c.last_updated_timestamp">never</span>
								<template v-else>
									<span class="badge text-bg-success" v-if="c.last_updated_success"><i class="fa fa-check"></i> ok</span>
									<span class="badge text-bg-danger" v-else><i class="fa fa-xmark"></i> failed</span>
									<span class="text-muted-bx ms-1">{{ format_time(c.last_updated_timestamp) }}</span>
									<div class="text-muted-bx text-break" v-if="c.last_updated_message">{{ c.last_updated_message }}</div>
								</template>
							</td>
							<td class="text-end text-nowrap">
								<button type="button" class="btn btn-sm btn-outline-success me-1" @click="run(c)" :disabled="running.has(c.id)" :title="'Fetch media from ' + c.hostname + ' now'">
									<span class="spinner-border spinner-border-sm" v-if="running.has(c.id)"></span>
									<i class="fa fa-play" v-else></i>
									Run
								</button>
								<button type="button" class="btn btn-sm btn-outline-primary me-1" @click="edit(c)" :title="'Edit ' + c.hostname">
									<i class="fa fa-pen"></i>
									<span class="visually-hidden">Edit</span>
								</button>
								<button type="button" class="btn btn-sm btn-outline-danger" @click="remove(c)" :title="'Delete ' + c.hostname">
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
