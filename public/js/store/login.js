
const store = Vue.reactive({
	loggedIn: false,
	/* { username, user_id, type, mail, permissions } */
	claims: null,
	token: null
});

export default store;

export function hasPermission(permission){
	return store.claims && store.claims.permissions.find(p => p == permission);
}