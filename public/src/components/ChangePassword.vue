<script setup lang="ts">
import { ref } from 'vue';
import { change_password as api_change_password } from "../api/user";

const props = defineProps<{
	useruid: string;
}>();

const old_password = ref("");
const new_password = ref("");
const password_wrong = ref(false);
const success = ref(false);
const busy = ref(false);

async function change_password() {
	password_wrong.value = false;
	success.value = false;
	busy.value = true;
	try {
		const s = await api_change_password(props.useruid, {
			old_password: old_password.value,
			new_password: new_password.value
		});
		password_wrong.value = !s;
		success.value = s;
		if (s) {
			old_password.value = "";
			new_password.value = "";
		}
	} finally {
		busy.value = false;
	}
}
</script>

<template>
	<form class="row g-2 align-items-start" @submit.prevent="change_password">
		<div class="col-sm-4">
			<label class="form-label small" for="pw-old">Current password</label>
			<input id="pw-old" type="password" v-model="old_password" autocomplete="current-password" class="form-control" :class="{'is-invalid': password_wrong}"/>
			<div class="invalid-feedback">The current password is wrong.</div>
		</div>
		<div class="col-sm-4">
			<label class="form-label small" for="pw-new">New password</label>
			<input id="pw-new" type="password" v-model="new_password" autocomplete="new-password" class="form-control"/>
		</div>
		<div class="col-sm-4">
			<label class="form-label small d-none d-sm-block">&nbsp;</label>
			<button class="btn btn-primary w-100" type="submit" :disabled="!old_password || !new_password || busy">
				<span class="spinner-border spinner-border-sm" v-if="busy"></span>
				Change password
			</button>
		</div>
	</form>
	<div class="alert alert-success py-2 small mt-3 mb-0" v-if="success">
		<i class="fa fa-check"></i> Your password has been changed.
	</div>
</template>
