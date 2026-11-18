<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { get_access_tokens, create_access_token, delete_access_token } from "../api/access_token";
import type { AccessToken } from "../api/types";

import ClipboardCopy from "./ClipboardCopy.vue";

defineProps<{
	username: string | null;
}>();

const list = ref<AccessToken[]>([]);
const new_name = ref("");
const new_days = ref(30);
const busy = ref(false);

const is_expired = (at: AccessToken) => Date.now() > at.expires;
const format_date = (ts: number) => new Date(ts).toLocaleDateString();

async function update() {
	list.value = await get_access_tokens();
}

async function create() {
	busy.value = true;
	try {
		await create_access_token({
			name: new_name.value,
			expires: Date.now() + (1000 * 3600 * 24 * new_days.value)
		});
		new_name.value = "";
		await update();
	} finally {
		busy.value = false;
	}
}

async function remove(at: AccessToken) {
	if (!confirm(`Delete the access token "${at.name}"? Servers using it will be logged out.`)) {
		return;
	}
	await delete_access_token(at.uid);
	await update();
}

onMounted(update);
</script>

<template>
	<p class="text-muted-bx">
		An access token lets the in-game mod upload in your name. Create one per server,
		then paste its <b>/bx_login</b> command into the chat on that server.
		Treat tokens like passwords and only use long-lived ones on servers you trust.
	</p>

	<form class="row g-2 align-items-end mb-4" @submit.prevent="create">
		<div class="col-sm-6">
			<label class="form-label small" for="token-name">Name</label>
			<input id="token-name" class="form-control" placeholder="e.g. the server name" v-model="new_name">
		</div>
		<div class="col-sm-3">
			<label class="form-label small" for="token-days">Valid for</label>
			<select id="token-days" class="form-select" v-model.number="new_days">
				<option :value="1">1 day</option>
				<option :value="7">7 days</option>
				<option :value="30">30 days</option>
				<option :value="90">90 days</option>
				<option :value="365">1 year</option>
			</select>
		</div>
		<div class="col-sm-3">
			<button class="btn btn-success w-100" type="submit" :disabled="!new_name || busy">
				<span class="spinner-border spinner-border-sm" v-if="busy"></span>
				<i class="fa fa-plus" v-else></i>
				Create token
			</button>
		</div>
	</form>

	<p class="text-muted-bx small mb-0" v-if="list.length == 0">You don't have any access tokens yet.</p>
	<div class="table-responsive" v-else>
		<table class="table table-hover align-middle mb-0">
			<thead>
				<tr>
					<th>Name</th>
					<th>Login command</th>
					<th class="d-none d-md-table-cell">Expires</th>
					<th class="d-none d-md-table-cell text-end">Used</th>
					<th></th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="at in list" :key="at.uid" :class="{'opacity-50': is_expired(at)}">
					<td>
						{{ at.name }}
						<span class="badge bg-danger ms-1" v-if="is_expired(at)">expired</span>
					</td>
					<td style="max-width: 22rem;">
						<clipboard-copy :text="'/bx_login ' + username + ' ' + at.token"/>
					</td>
					<td class="d-none d-md-table-cell small">{{ format_date(at.expires) }}</td>
					<td class="d-none d-md-table-cell small text-end">{{ at.usecount }}×</td>
					<td class="text-end">
						<button type="button" class="btn btn-sm btn-outline-danger" @click="remove(at)" :title="'Delete ' + at.name">
							<i class="fa fa-trash"></i>
							<span class="visually-hidden">Delete</span>
						</button>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>
