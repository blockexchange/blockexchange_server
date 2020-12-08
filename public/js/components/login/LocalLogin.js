import store from '../../store/login.js';
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
		logout: LoginService.logout,
		isLoggedIn: function() {
			return store.loggedIn;
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
				Login <span class="badge badge-danger">{{ message }}</span>
			</button>
			<button v-if="isLoggedIn()"
				class="btn btn-secondary btn-block"
				v-on:click="logout()">
				Logout
			</button>
		</form>
	`
};
