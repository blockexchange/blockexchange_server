import store from '../../store/login.js';

export default {
	data: function(){
		return {
			store: store
		};
	},
	template: /*html*/`
		<div v-if="store.loggedIn">
			<router-link to="/profile">
				<i class="fas fa-user"></i>
				<span>
					Logged in as <b>{{ store.claims.username }}</b>
				</span>
			</router-link>
		</div>
	`
};
