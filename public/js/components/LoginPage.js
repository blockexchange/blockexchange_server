import { login } from '../service/login.js';
import store from '../store/token.js';

export default {
	data: function() {
		return {
			username: "",
			password: "",
			store: store
		};
	},
	methods: {
		login: login
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
					<button v-if="!store.token" class="btn btn-secondary btn-block" v-on:click="login(username, password)">
						Login
					</button>
				</form>
			</div>
			<div class="col-md-4"></div>
		</div>
	`
};
