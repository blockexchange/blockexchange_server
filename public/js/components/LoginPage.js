import { request_token } from '../api/token.js';
import store from '../store/token.js';

export default {
	data: function() {
		return {
			username: "",
			password: "",
			message: null
		};
	},
	methods: {
		login(username, password) {
			this.message = null;
		  request_token(username, password)
		  .then(t => store.token = t)
		  .catch(e => this.message = e.message);
		},
		logout() {
			store.token = null;
		},
		isLoggedIn() {
			return !!store.token;
		}
	},
	template: /*html*/`
		<div class="row">
			<div class="col-md-4"></div>
			<div class="col-md-4">
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
					<button v-if="!isLoggedIn()" class="btn btn-secondary btn-block" v-on:click="login(username, password)">
						Login
					</button>
					<button v-if="isLoggedIn()" class="btn btn-secondary btn-block" v-on:click="logout()">
						Logout
					</button>
					<span v-if="message">
						{{ message }}
					</span>
				</form>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
