<script setup lang="ts">
import { ref } from 'vue';
import { unlink_oauth } from "../api/user";
import { check_login } from "../service/login";

const props = defineProps<{
	useruid: string;
}>();

const password_1 = ref("");
const password_2 = ref("");
const busy = ref(false);

async function unlink() {
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
	<h5>Unlink OAuth provider</h5>
	<div class="input-group">
		<input type="password" v-model="password_1" placeholder="New password" class="form-control"/>
		<input type="password" v-model="password_2" placeholder="New password (repeat)" class="form-control"/>
		<button class="btn btn-warning" :disabled="!password_1 || !password_2 || password_1 != password_2 || busy" @click="unlink">
			<i class="fa fa-link-slash" v-if="!busy"></i>
			<i class="fa fa-spinner fa-spin" v-if="busy"></i>
			Unlink and set password
		</button>
	</div>
	<div class="alert alert-warning">
		<i class="fa fa-triangle-exclamation"></i>
		<b>Warning:</b> This unlinks the account from the OAuth provider, this action is irreversible!
	</div>
</template>
