import store from '../../store/login.js';
import debounce from '../../util/debounce.js';
import { validate_username, update } from '../../api/user.js';
import loginservice from '../../service/login.js';

export default {
	data: function(){
		return {
			store: store,
			username: store.claims.username,
			mail: store.claims.mail,
			username_valid: true,
			username_message: null
		};
	},
	computed: {
		enable_update: function(){
			const values_changed = (
				this.username !== store.claims.username ||
				this.mail !== store.claims.mail
			);
			return values_changed && this.username_valid;
		}
	},
	watch: {
		username: debounce(function(newUsername){
			if (newUsername === store.claims.username){
				// same as before, valid
				this.username_valid = true;
				this.username_message = null;
				return;
			}

			validate_username(newUsername)
			.then(result => {
				this.username_valid = result.valid;
				this.username_message = result.message;
			});
		}, 250)
	},
	methods: {
		update: function(){
			update({
				id: store.claims.user_id,
				name: this.username,
				mail: this.mail
			})
			.then(() => {
				store.claims.mail = this.mail;
				store.claims.username = this.username;
				loginservice.persist();
			});
		}
	},
	template: /*html*/`
		<div>
			<input type="text"
				class="form-control"
				placeholder="Username"
				v-bind:class="{ 'is-invalid': !username_valid }"
				v-model="username"
			/>
			<div class="alert alert-danger" v-if="username_message">
				<b>Error:</b> {{ username_message }}
			</div>
			<div class="alert alert-warning" v-if="username !== store.claims.username">
				<b>Warning:</b> if you change your username all your schemas will be only available
				under the new name, this will affect all your bookmarks too!
			</div>
			<input type="text"
				class="form-control"
				placeholder="E-Mail"
				v-model="mail"
			/>
			<div class="alert alert-success">
				E-Mail address is optional but strongly encouraged as it is the only recovery-option
				for local accounts
			</div>
			<button class="btn btn-primary btn-block" v-bind:disabled="!enable_update" v-on:click="update">
				<i class="fa fa-save"/> Update profile
			</button>
		</div>
	`
};
