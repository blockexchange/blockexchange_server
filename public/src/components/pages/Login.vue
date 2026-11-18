<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { login as service_login, logout as service_logout, get_username } from "../../service/login";
import { get_cdb_login, get_discord_login, get_github_login, get_mesehub_login, get_codeberg_login, get_base_url } from "../../service/info";

const router = useRouter();

const username = ref("");
const password = ref("");
const busy = ref(false);
const login_failed = ref(false);

const current_user = computed(get_username);
const valid_input = computed(() => username.value != "" && password.value != "");

const base_url = computed(get_base_url);
const providers = computed(() => [
	{ name: "GitHub", url: get_github_login(), icon: "fa-brands fa-github" },
	{ name: "Discord", url: get_discord_login(), icon: "fa-brands fa-discord" },
	{ name: "ContentDB", url: get_cdb_login(), img: "/pics/contentdb.png" },
	{ name: "Codeberg", url: get_codeberg_login(), img: "/pics/codeberg.png" },
	{ name: "Mesehub", url: get_mesehub_login(), img: "/pics/default_mese_crystal.png" }
].filter(p => p.url));

async function login() {
	busy.value = true;
	login_failed.value = false;
	try {
		const claims = await service_login(username.value, password.value);
		if (!claims) {
			login_failed.value = true;
			return;
		}
	} finally {
		busy.value = false;
	}
	await router.push("/user/" + get_username());
}

async function logout() {
	busy.value = true;
	try {
		await service_logout();
	} finally {
		busy.value = false;
	}
}
</script>

<template>
	<div class="row justify-content-center">
		<div class="col-sm-10 col-md-7 col-lg-5">
			<div class="card">
				<div class="card-body p-4">
					<template v-if="current_user">
						<h1 class="h4 mb-3">You're logged in</h1>
						<p class="text-muted-bx">Signed in as <b class="text-light">{{ current_user }}</b>.</p>
						<div class="d-grid gap-2">
							<router-link :to="'/user/' + current_user" class="btn btn-success">Go to my schematics</router-link>
							<button class="btn btn-outline-light" type="button" :disabled="busy" @click="logout">
								<span class="spinner-border spinner-border-sm" v-if="busy"></span>
								Log out
							</button>
						</div>
					</template>
					<template v-else>
						<h1 class="h4 mb-1">Log in</h1>
						<p class="text-muted-bx small mb-4">Log in to upload, star and manage your schematics.</p>

						<div class="d-grid gap-2" v-if="providers.length">
							<a v-for="p in providers" :key="p.name" :href="p.url" class="btn btn-secondary bx-provider-btn">
								<i :class="p.icon" v-if="p.icon"></i>
								<img :src="(base_url ?? '') + (p.img ?? '')" alt="" v-else>
								Continue with {{ p.name }}
							</a>
						</div>
						<div class="bx-divider" v-if="providers.length">or with a username</div>

						<form @submit.prevent="login" novalidate>
							<div class="alert alert-danger py-2 small" v-if="login_failed" role="alert">
								<i class="fa fa-circle-exclamation"></i>
								Wrong username or password.
							</div>
							<div class="mb-3">
								<label class="form-label" for="login-username">Username</label>
								<input id="login-username" type="text" class="form-control" autocomplete="username" autofocus v-model="username"/>
							</div>
							<div class="mb-3">
								<label class="form-label" for="login-password">Password</label>
								<input id="login-password" type="password" class="form-control" autocomplete="current-password" v-model="password"/>
							</div>
							<button class="btn btn-success w-100" type="submit" :disabled="!valid_input || busy">
								<span class="spinner-border spinner-border-sm" v-if="busy"></span>
								Log in
							</button>
						</form>
					</template>
				</div>
			</div>
			<p class="text-center small text-muted-bx mt-3" v-if="!current_user">
				New here? <router-link to="/register">Create an account</router-link>
			</p>
		</div>
	</div>
</template>
