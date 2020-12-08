export default {
	beforeRouteEnter: function(to, from, next){
		console.log(to.params.token);
		next({ path: "/login" });
	},
	template: `<div></div>`
};
