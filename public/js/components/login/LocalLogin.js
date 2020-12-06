import store from '../../store/token.js';
import LoginService from '../../service/login.js';

export default {
	data: function() {
		return {
			username: "",
			password: "",
			message: null
		};
	},
	methods: {
		login: function(username, password) {
			this.message = null;
			LoginService.login(username, password)
			.then(result => {
				if (!result.success){
					this.message = result.message;
				}
			});
		},
		isLoggedIn: function() {
			return !!store.token;
		}
	},
	template: /*html*/`
		<form v-on:submit.prevent>
			<input type="text"
				class="form-control"
				placeholder="Username"
				v-model="username"
			/>
			<input type="password"
				class="form-control"
				placeholder="Password"
				v-model="password"
			/>
			<button v-if="!isLoggedIn()"
				v-bind:disabled="!username || !password"
				class="btn btn-secondary btn-block"
				v-on:click="login(username, password)">
				Login
			</button>
			<button v-if="isLoggedIn()"
				class="btn btn-secondary btn-block"
				v-on:click="loginService.logout()">
				Logout
			</button>
			<span v-if="message">
				<div class="alert alert-danger" role="alert">
				  {{ message }}
				</div>
			</span>
		</form>
	`
};
