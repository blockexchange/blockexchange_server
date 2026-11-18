<script setup lang="ts">
import { onMounted, ref } from 'vue';

import { get_access_tokens, create_access_token, delete_access_token } from "../api/access_token";
import type { AccessToken } from "../api/types";
import format_time from "../util/format_time";

import ClipboardCopy from "./ClipboardCopy.vue";

defineProps<{
	username: string | null;
}>();

const list = ref<AccessToken[]>([]);
const new_name = ref("");
const new_days = ref(365);

const is_expired = (at: AccessToken) => Date.now() > at.expires;

async function update() {
	list.value = await get_access_tokens();
}

async function create() {
	await create_access_token({
		name: new_name.value,
		expires: Date.now() + (1000 * 3600 * 24 * new_days.value)
	});
	new_name.value = "";
	await update();
}

async function remove(at: AccessToken) {
	await delete_access_token(at.uid);
	await update();
}

onMounted(update);
</script>

<template>
	<div>
		<h5>Manage access-tokens</h5>
		<div class="alert alert-secondary">
			<i class="fa fa-circle-info"></i>
			<b>Note:</b> Access-tokens serve as a time-limited login mechanism.
			Tokens with long expiration-times should only be used on trusted servers.
		</div>
		<table class="table table-striped table-condensed table-dark">
			<thead>
				<tr>
					<th>Name</th>
					<th>Token/Login</th>
					<th>Created</th>
					<th>Expires</th>
					<th>Use-count</th>
					<th>Action</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="at in list" :key="at.uid" :class="{'table-warning': is_expired(at)}">
					<td>{{ at.name }}</td>
					<td>
						<clipboard-copy :text="'/bx_login ' + username + ' ' + at.token"/>
					</td>
					<td>{{ format_time(at.created) }}</td>
					<td>{{ format_time(at.expires) }}</td>
					<td>{{ at.usecount }}</td>
					<td>
						<a class="btn btn-danger" @click="remove(at)">
							<i class="fa fa-trash"></i>
							Delete
						</a>
					</td>
				</tr>
				<tr>
					<td colspan="2">
						<label>Token-Name (usually the server-name you are using it on)</label>
						<input class="form-control" placeholder="Name" v-model="new_name">
					</td>
					<td></td>
					<td>
						<label>Expiration (in days)</label>
						<input class="form-control" type="number" placeholder="Expiration (days)" v-model.number="new_days">
					</td>
					<td></td>
					<td>
						<a class="btn btn-primary" :class="{disabled: !new_name || !new_days || new_days <= 0}" @click="create">
							<i class="fa fa-plus"></i>
							Create
						</a>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>
