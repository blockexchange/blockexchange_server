import Breadcrumb, { START, REGISTER } from "../Breadcrumb.js";
import { create_captcha } from "../../api/captcha.js";
import { register } from "../../api/register.js";
import { login } from "../../service/login.js";

export default {
	components: {
        "bread-crumb": Breadcrumb
	},
	mounted: function() {
		this.update_captcha();
	},
	data: function() {
		return {
			name: "",
			password: "",
			password2: "",
			captcha_answer: "",
			captcha_id: null,
			register_result: {},
			breadcrumb: [START, REGISTER]
		};
	},
	methods: {
		update_captcha: function() {
			create_captcha().then(c => this.captcha_id = c);
		},
		register: function() {
			const rr = {
				name: this.name,
				password: this.password,
				captcha_id: this.captcha_id,
				captcha_answer: this.captcha_answer
			};

			register(rr)
			.then(r => {
				this.register_result = r;
				if (!r.success) {
					// update captcha
					this.update_captcha();
				} else {
					// register and redirect
					login(this.name, this.password)
					.then(() => this.$router.push("/profile"));
				}
			});
		}
	},
	computed: {
		captcha_src: function() {
			return this.captcha_id ? `${BaseURL}/api/captcha/${this.captcha_id}.png` : null;
		},
		can_register: function() {
			return this.name && this.password && this.password == this.password2 && this.captcha_answer;
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<div class="row">
            <div class="col-md-4"></div>
            <div class="col-md-4 card" style="padding: 20px;">
				<form @submit.prevent>
					<h5>Register a new account</h5>
					<div class="input-group has-validation">
						<input type="text"
							v-model="name" class="form-control w-100"
							v-bind:class="{'is-invalid': register_result.error_username_taken || register_result.error_invalid_username}"
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
							v-bind:class="{'is-invalid': register_result.error_password_too_short}"
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
							v-bind:class="{'is-invalid': register_result.error_captcha}"
							placeholder="Captcha answer"/>
						<div class="invalid-feedback" v-if="register_result.error_captcha">
							Captcha answer invalid
						</div>
					</div>
					<button type="submit" class="btn btn-primary w-100" :disabled="!can_register" v-on:click="register">
						<i class="fa fa-user-plus"></i>
						Register
					</button>
				</form>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
