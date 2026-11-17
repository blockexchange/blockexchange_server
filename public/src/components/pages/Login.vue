<script setup lang="ts">
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

import { login as service_login, logout as service_logout, is_logged_in } from "../../service/login";
import { get_cdb_login, get_discord_login, get_github_login, get_mesehub_login, get_codeberg_login, get_base_url } from "../../service/info";
import BreadCrumb from "../Breadcrumb.vue";
import { LOGIN, START } from "../../util/breadcrumbs";

const router = useRouter();

const breadcrumb = [START, LOGIN];
const username = ref("");
const password = ref("");
const busy = ref(false);
const error_message = ref("");

const logged_in = computed(is_logged_in);
const valid_input = computed(() => username.value != "" && password.value != "");

const github_login = computed(get_github_login);
const cdb_login = computed(get_cdb_login);
const discord_login = computed(get_discord_login);
const mesehub_login = computed(get_mesehub_login);
const codeberg_login = computed(get_codeberg_login);
const base_url = computed(get_base_url);

function login() {
	busy.value = true;
	error_message.value = "";
	service_login(username.value, password.value)
	.then(claims => {
		busy.value = false;
		if (!claims) {
			// no luck
			error_message.value = "Login failed!";
		} else {
			// go to base page
			router.push("/profile");
		}
	});
}

function logout() {
	busy.value = true;
	service_logout()
	.then(() => busy.value = false);
}
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="row">
		<div class="col-md-2"></div>
		<div class="col-md-8 card" style="padding: 20px;">
			<div class="row">
				<div class="col-md-6">
					<h5>Login with username and password</h5>
					<form @submit.prevent="login">
						<input type="text"
							class="form-control"
							placeholder="Username"
							:disabled="logged_in"
							v-model="username"/>
						<input type="password"
							class="form-control"
							placeholder="Password"
							:disabled="logged_in"
							v-model="password"/>
						<button class="btn btn-primary w-100" v-if="!logged_in" type="submit" :disabled="!valid_input">
							<i class="fa-solid fa-right-to-bracket"></i>
							Login
							<i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
							<span class="badge bg-danger">{{ error_message }}</span>
						</button>
						<a class="btn btn-secondary w-100" v-if="logged_in" @click="logout">
							<i class="fa-solid fa-right-from-bracket"></i>
							Logout
							<i class="fa-solid fa-spinner fa-spin" v-if="busy"></i>
						</a>
						Register a new account <router-link to="/register">here</router-link>
					</form>
				</div>
				<div class="col-md-6">
					<h5>Login with external provider</h5>
					<a :href="github_login" class="btn btn-secondary w-100" :class="{disabled: logged_in}" v-if="github_login">
						<i class="fab fa-github"></i>
						Login with Github
					</a>
					&nbsp;
					<a :href="cdb_login" class="btn btn-secondary w-100" :class="{disabled: logged_in}" v-if="cdb_login">
						<img :src="base_url + '/pics/contentdb.png'" height="24" width="24">
						Login with ContentDB
					</a>
					&nbsp;
					<a :href="discord_login" class="btn btn-secondary w-100" :class="{disabled: logged_in}" v-if="discord_login">
						<i class="fab fa-discord"></i>
						Login with Discord
					</a>
					&nbsp;
					<a :href="mesehub_login" class="btn btn-secondary w-100" :class="{disabled: logged_in}" v-if="mesehub_login">
						<img :src="base_url + '/pics/default_mese_crystal.png'">
						Login with Mesehub
					</a>
					&nbsp;
					<a :href="codeberg_login" class="btn btn-secondary w-100" :class="{disabled: logged_in}" v-if="codeberg_login">
						<img :src="base_url + '/pics/codeberg.png'" height="24" width="24">
						Login with Codeberg
					</a>
				</div>
			</div>
		</div>
		<div class="col-md-4"></div>
	</div>
</template>
