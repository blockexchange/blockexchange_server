<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { search_users, save_user, get_user } from "../api/user";
import { check_login, get_claims } from "../service/login";

const name = ref("");
const username_taken = ref(false);

const can_rename = computed(() => !!name.value && name.value != get_claims()?.username);
const username_valid = computed(() => /^[a-zA-Z0-9_.-]*$/.test(name.value));
const is_valid = computed(() => username_valid.value && !username_taken.value);

watch(name, () => username_taken.value = false);

onMounted(() => {
	name.value = get_claims()?.username ?? "";
});

async function rename() {
	const list = await search_users({ name: name.value });
	username_taken.value = list.length > 0;
	const claims = get_claims();
	if (username_taken.value || !claims) {
		return;
	}
	const user = await get_user(claims.user_uid);
	user.name = name.value;
	await save_user(user);
	await check_login(true);
	location.reload();
}
</script>

<template>
	<h5>Change username</h5>
	<div class="input-group has-validation">
		<input type="text" v-model="name" class="form-control" :class="{'is-invalid': !is_valid}"/>
		<button class="btn btn-warning" :disabled="!can_rename || !is_valid" @click="rename">
			<i class="fa fa-edit"></i>
			Rename
		</button>
		<div class="invalid-feedback" v-if="!username_valid">
			Username is invalid, allowed chars: a to z, A to Z, 0 to 9 and -, _
		</div>
		<div class="invalid-feedback" v-if="username_taken">
			Username is already taken
		</div>
	</div>
	<div class="alert alert-warning">
		<i class="fa fa-triangle-exclamation"></i>
		<b>Warning:</b> Changing the user name affects all your schematic urls
	</div>
</template>
