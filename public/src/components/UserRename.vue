<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { search_users, save_user, get_user } from "../api/user";
import { check_login, get_claims } from "../service/login";

const name = ref("");
const username_taken = ref(false);
const busy = ref(false);

const can_rename = computed(() => !!name.value && name.value != get_claims()?.username);
const username_valid = computed(() => /^[a-zA-Z0-9_.-]*$/.test(name.value));
const is_valid = computed(() => username_valid.value && !username_taken.value);

watch(name, () => username_taken.value = false);

onMounted(() => {
	name.value = get_claims()?.username ?? "";
});

async function rename() {
	const claims = get_claims();
	if (!claims || !confirm(`Rename "${claims.username}" to "${name.value}"? Existing links to your schematics will stop working.`)) {
		return;
	}
	busy.value = true;
	try {
		const list = await search_users({ name: name.value });
		username_taken.value = list.length > 0;
		if (username_taken.value) {
			return;
		}
		const user = await get_user(claims.user_uid);
		user.name = name.value;
		await save_user(user);
		await check_login(true);
		location.reload();
	} finally {
		busy.value = false;
	}
}
</script>

<template>
	<p class="text-muted-bx">
		Your username is part of every schematic link and in-game command, such as <code>/bx_load {{ name || 'name' }} …</code>.
		Renaming breaks links and commands that others have already shared.
	</p>
	<form class="row g-2 align-items-start" @submit.prevent="rename">
		<div class="col-sm-8">
			<label class="form-label small visually-hidden" for="rename-name">New username</label>
			<input id="rename-name" type="text" v-model="name" class="form-control" :class="{'is-invalid': !is_valid}"/>
			<div class="invalid-feedback" v-if="!username_valid">
				Only letters (a-z, A-Z), numbers and - _ . are allowed.
			</div>
			<div class="invalid-feedback" v-if="username_taken">
				This username is already taken.
			</div>
		</div>
		<div class="col-sm-4">
			<button class="btn btn-outline-warning w-100" type="submit" :disabled="!can_rename || !is_valid || busy">
				<span class="spinner-border spinner-border-sm" v-if="busy"></span>
				Rename account
			</button>
		</div>
	</form>
</template>
