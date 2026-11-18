<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import { create_captcha } from "../../api/captcha";
import { register as api_register } from "../../api/register";
import type { RegisterResponse } from "../../api/types";
import { login } from "../../service/login";
import { get_cdb_login, get_codeberg_login, get_discord_login, get_github_login, get_mesehub_login } from "../../service/info";

const router = useRouter();

const name = ref("");
const password = ref("");
const password2 = ref("");
const captcha_answer = ref("");
const captcha_id = ref<string | null>(null);
const register_result = ref<Partial<RegisterResponse>>({});
const busy = ref(false);
const has_providers = computed(() => !!(get_github_login() || get_discord_login() || get_cdb_login() || get_codeberg_login() || get_mesehub_login()));

const captcha_src = computed(() => captcha_id.value ? `/api/captcha/${captcha_id.value}.png` : null);
const name_valid = computed(() => /^[a-zA-Z0-9_-]*$/.test(name.value));
const password_mismatch = computed(() => password2.value != "" && password.value != password2.value);
const can_register = computed(() => name.value && name_valid.value && password.value && password.value == password2.value && captcha_answer.value);

async function update_captcha() {
	captcha_id.value = null;
	captcha_answer.value = "";
	captcha_id.value = await create_captcha();
}

async function register() {
	busy.value = true;
	try {
		const r = await api_register({
			name: name.value,
			password: password.value,
			captcha_id: captcha_id.value ?? "",
			captcha_answer: captcha_answer.value
		});
		register_result.value = r;
		if (!r.success) {
			// a captcha can only be used once
			await update_captcha();
			return;
		}
		await login(name.value, password.value);
	} finally {
		busy.value = false;
	}
	await router.push("/profile");
}

onMounted(update_captcha);
</script>

<template>
	<div class="row justify-content-center">
		<div class="col-sm-10 col-md-7 col-lg-5">
			<div class="card">
				<div class="card-body p-4">
					<h1 class="h4 mb-1">Create an account</h1>
					<p class="text-muted-bx small mb-4">
						An account lets you upload schematics and star your favourites.
						<template v-if="has_providers">
							You can also <router-link to="/login">log in with an existing account</router-link> (e.g. GitHub) without signing up.
						</template>
					</p>
					<form @submit.prevent="register" novalidate>
						<div class="mb-3">
							<label class="form-label" for="reg-name">Username</label>
							<input id="reg-name" type="text" v-model="name" class="form-control" autocomplete="username"
								:class="{'is-invalid': !name_valid || register_result.error_username_taken || register_result.error_invalid_username}"/>
							<div class="form-text" v-if="name_valid && !register_result.error_username_taken && !register_result.error_invalid_username">
								Letters, numbers, - and _. This is also used in your schematic links.
							</div>
							<div class="invalid-feedback" v-if="register_result.error_username_taken">
								This username is already taken.
							</div>
							<div class="invalid-feedback" v-else-if="!name_valid || register_result.error_invalid_username">
								Only letters (a-z, A-Z), numbers, - and _ are allowed.
							</div>
						</div>
						<div class="mb-3">
							<label class="form-label" for="reg-password">Password</label>
							<input id="reg-password" type="password" v-model="password" class="form-control" autocomplete="new-password"
								:class="{'is-invalid': register_result.error_password_too_short}"/>
							<div class="form-text" v-if="!register_result.error_password_too_short">At least 6 characters.</div>
							<div class="invalid-feedback">The password needs at least 6 characters.</div>
						</div>
						<div class="mb-3">
							<label class="form-label" for="reg-password2">Repeat password</label>
							<input id="reg-password2" type="password" v-model="password2" class="form-control" autocomplete="new-password"
								:class="{'is-invalid': password_mismatch}"/>
							<div class="invalid-feedback">The passwords don't match.</div>
						</div>
						<div class="mb-4">
							<label class="form-label" for="reg-captcha">Type the characters you see</label>
							<div class="d-flex align-items-center gap-2 mb-2">
								<img :src="captcha_src" v-if="captcha_src" alt="Captcha" class="rounded bg-light">
								<span class="text-muted-bx small" v-else>
									<span class="spinner-border spinner-border-sm"></span> Loading captcha…
								</span>
								<button type="button" class="btn btn-sm btn-outline-secondary" @click="update_captcha" title="Show a different captcha">
									<i class="fa fa-rotate"></i>
									<span class="visually-hidden">New captcha</span>
								</button>
							</div>
							<input id="reg-captcha" type="text" v-model="captcha_answer" class="form-control" autocomplete="off"
								:class="{'is-invalid': register_result.error_captcha}"/>
							<div class="invalid-feedback">That wasn't right, please try the new captcha.</div>
						</div>
						<button type="submit" class="btn btn-success w-100" :disabled="!can_register || busy">
							<span class="spinner-border spinner-border-sm" v-if="busy"></span>
							Create account
						</button>
					</form>
				</div>
			</div>
			<p class="text-center small text-muted-bx mt-3">
				Already have an account? <router-link to="/login">Log in</router-link>
			</p>
		</div>
	</div>
</template>
