import store from '../../store/login.js';

export default {
	data: function(){
		return {
			store: store
		};
	},
	template: /*html*/`
		<span v-if="store.loggedIn" class="badge badge-success">Logged in as {{ store.username }}</span>
	`
};
