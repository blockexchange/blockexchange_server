import loginStore from '../../store/login.js';

export default {
	beforeRouteEnter: function(to, from, next){
		const payload = JSON.parse(atob(to.params.token.split(".")[1]));
		loginStore.token = to.params.token;
		loginStore.loggedIn = true;
		loginStore.username = payload.username;

		next({ path: "/login" });
	},
	template: `<div></div>`
};
