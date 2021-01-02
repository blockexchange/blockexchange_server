import debounce from '../../util/debounce.js';
import { validate_username } from '../../api/user.js';
import { register } from '../../api/register.js';
import loginService from '../../service/login.js';

export default {
	data: function(){
		return {
			username: "",
			password: "",
			password2: "",
			mail: "",
			username_valid: true,
			username_message: "",
			response_message: ""
		};
	},
	methods: {
		register: function(){
			this.response_message = "";
			register(this.username, this.password)
			.then(r => {
				if (r.success){
					// signup succeeded
					return loginService.login(this.username, this.password);
				} else {
					// signup failed
					this.response_message = r.message;
				}
			})
			.then(r => {
				if (r && r.success){
					// signup and login ok
					this.$router.push("/profile");
				}
			});
		}
	},
	watch: {
		username: debounce(function(newUsername){
			validate_username(newUsername)
			.then(result => {
				this.username_valid = result.valid;
				this.username_message = result.message;
			});
		}, 250)
	},
	computed: {
		register_enabled: function(){
			return this.username &&
				this.password &&
				this.username_valid &&
				this.password == this.password2;
		},
		password_mismatch: function(){
			return this.password &&
				this.password != this.password2;
		}
	},
	template: /*html*/`
		<form v-on:submit.prevent>
			<input type="text"
				class="form-control"
				placeholder="Username"
				v-model="username"
				v-bind:class="{ 'is-invalid': !username_valid, 'is-valid': username && username_valid }"
			/>
			<div class="alert alert-danger" v-if="!username_valid">
				<b>Error:</b> {{ username_message }}
			</div>
			<input type="text"
				class="form-control"
				placeholder="E-Mail (optional)"
				v-model="mail"
			/>
			<input type="password"
				class="form-control"
				placeholder="Password"
				v-model="password"
				v-bind:class="{ 'is-invalid': password_mismatch }"
			/>
			<input type="password"
				class="form-control"
				placeholder="Password (verify)"
				v-model="password2"
				v-bind:class="{ 'is-invalid': password_mismatch }"
			/>
			<button class="btn btn-secondary btn-block" v-on:click="register" :disabled="!register_enabled">
				Register
			</button>
			<div class="alert alert-danger" v-if="response_message">
				<b>Error:</b> {{ response_message }}
			</div>
		</form>
	`
};
