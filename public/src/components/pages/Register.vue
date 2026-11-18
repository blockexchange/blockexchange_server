<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';

import BreadCrumb from "../Breadcrumb.vue";
import { START, REGISTER } from "../../util/breadcrumbs";
import { create_captcha } from "../../api/captcha";
import { register as api_register } from "../../api/register";
import type { RegisterResponse } from "../../api/types";
import { login } from "../../service/login";

const router = useRouter();

const breadcrumb = [START, REGISTER];
const name = ref("");
const password = ref("");
const password2 = ref("");
const captcha_answer = ref("");
const captcha_id = ref<string | null>(null);
const register_result = ref<Partial<RegisterResponse>>({});

const captcha_src = computed(() => captcha_id.value ? `/api/captcha/${captcha_id.value}.png` : null);
const can_register = computed(() => name.value && password.value && password.value == password2.value && captcha_answer.value);

async function update_captcha() {
	captcha_id.value = await create_captcha();
}

async function register() {
	const r = await api_register({
		name: name.value,
		password: password.value,
		captcha_id: captcha_id.value ?? "",
		captcha_answer: captcha_answer.value
	});
	register_result.value = r;
	if (!r.success) {
		// update captcha
		await update_captcha();
	} else {
		// register and redirect
		await login(name.value, password.value);
		await router.push("/profile");
	}
}

onMounted(update_captcha);
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<div class="row">
		<div class="col-md-4"></div>
		<div class="col-md-4 card" style="padding: 20px;">
			<form @submit.prevent>
				<h5>Register a new account</h5>
				<div class="input-group has-validation">
					<input type="text"
						v-model="name" class="form-control w-100"
						:class="{'is-invalid': register_result.error_username_taken || register_result.error_invalid_username}"
						placeholder="Username"/>
					<div class="invalid-feedback" v-if="register_result.error_username_taken">
						Username is already taken
					</div>
					<div class="invalid-feedback" v-if="register_result.error_invalid_username">
						Username is invalid, allowed chars: a to z, A to Z, 0 to 9 and -, _
					</div>
				</div>
				<div class="input-group has-validation">
					<input type="password" v-model="password" class="form-control w-100"
						:class="{'is-invalid': register_result.error_password_too_short}"
						placeholder="Password"/>
					<div class="invalid-feedback" v-if="register_result.error_password_too_short">
						Password-length must be 6 characters or more
					</div>
				</div>
				<input type="password" v-model="password2" class="form-control w-100" placeholder="Retype password"/>
				<img :src="captcha_src" v-if="captcha_src">
				<div class="alert alert-secondary" v-if="!captcha_src">
					<i class="fa fa-spinner fa-spin"></i>
					Loading captcha
				</div>
				<div class="input-group has-validation">
					<input type="text" v-model="captcha_answer" class="form-control w-100"
						:class="{'is-invalid': register_result.error_captcha}"
						placeholder="Captcha answer"/>
					<div class="invalid-feedback" v-if="register_result.error_captcha">
						Captcha answer invalid
					</div>
				</div>
				<button type="submit" class="btn btn-primary w-100" :disabled="!can_register" @click="register">
					<i class="fa fa-user-plus"></i>
					Register
				</button>
			</form>
		</div>
		<div class="col-md-4"></div>
	</div>
</template>
