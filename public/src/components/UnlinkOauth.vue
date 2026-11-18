<script setup lang="ts">
import { computed, ref } from 'vue';
import { unlink_oauth } from "../api/user";
import { check_login } from "../service/login";

const props = defineProps<{
	useruid: string;
	provider: string;
}>();

const password_1 = ref("");
const password_2 = ref("");
const busy = ref(false);

const mismatch = computed(() => password_2.value != "" && password_1.value != password_2.value);

async function unlink() {
	if (!confirm(`Unlink your account from ${props.provider}? This can't be undone.`)) {
		return;
	}
	busy.value = true;
	try {
		await unlink_oauth(props.useruid, {
			new_password: password_1.value
		});
		await check_login(true);
		location.reload();
	} finally {
		busy.value = false;
	}
}
</script>

<template>
	<p class="text-muted-bx">
		You log in with <b class="text-body-emphasis">{{ provider.toLowerCase() }}</b>.
		To log in with a username and password instead, set a password below.
		This permanently removes the link to {{ provider.toLowerCase() }}.
	</p>
	<form class="row g-2 align-items-start" @submit.prevent="unlink">
		<div class="col-sm-4">
			<label class="form-label small" for="unlink-pw1">New password</label>
			<input id="unlink-pw1" type="password" v-model="password_1" autocomplete="new-password" class="form-control"/>
		</div>
		<div class="col-sm-4">
			<label class="form-label small" for="unlink-pw2">Repeat password</label>
			<input id="unlink-pw2" type="password" v-model="password_2" autocomplete="new-password" class="form-control" :class="{'is-invalid': mismatch}"/>
			<div class="invalid-feedback">The passwords don't match.</div>
		</div>
		<div class="col-sm-4">
			<label class="form-label small d-none d-sm-block">&nbsp;</label>
			<button class="btn btn-warning w-100" type="submit" :disabled="!password_1 || !password_2 || password_1 != password_2 || busy">
				<span class="spinner-border spinner-border-sm" v-if="busy"></span>
				<i class="fa fa-link-slash" v-else></i>
				Unlink and set password
			</button>
		</div>
	</form>
</template>
