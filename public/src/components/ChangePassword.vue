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
	} finally {
		busy.value = false;
	}
}
</script>

<template>
	<h5>Change password</h5>
	<div class="input-group has-validation">
		<input type="password" v-model="old_password" placeholder="Old password" class="form-control" :class="{'is-invalid': password_wrong}"/>
		<input type="password" v-model="new_password" placeholder="New password" class="form-control" :class="{'is-valid': success}"/>
		<button class="btn btn-primary" :disabled="!old_password || !new_password || busy" @click="change_password">
			<i class="fa fa-edit" v-if="!busy"></i>
			<i class="fa fa-spinner fa-spin" v-if="busy"></i>
			Change password
		</button>
		<div class="invalid-feedback" v-if="password_wrong">
			Password invalid
		</div>
		<div class="valid-feedback" v-if="success">
			Password changed
		</div>
	</div>
</template>
