import loginService from '../../service/login.js';

export default {
	beforeRouteEnter: function(to, from, next){
		loginService.parse_token(to.params.token);
		next({ path: "/profile" });
	},
	template: `<div></div>`
};