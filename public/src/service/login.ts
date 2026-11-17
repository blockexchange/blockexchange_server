import { reactive } from 'vue';

import { get_claims as fetch_claims, login as api_login, logout as api_logout } from '../api/login';
import type { Claims, Permission } from '../api/types';

const store = reactive<{ claims: Claims | null }>({
	claims: null
});

export const is_logged_in = () => store.claims != null;
export const get_claims = () => store.claims;

export const get_user_uid = () => store.claims ? store.claims.user_uid : null;
export const get_user_type = () => store.claims ? store.claims.type : null;
export const get_username = () => store.claims ? store.claims.username : null;
export const has_permission = (permission: Permission) => store.claims ? store.claims.permissions.includes(permission) : false;

export const check_login = (renew?: boolean) => fetch_claims(renew).then(c => {
	store.claims = c;
	return c;
});

export const login = (name: string, password: string) => {
	return api_login(name, password)
	.then(success => {
		if (success) {
			return check_login();
		}
		return null;
	});
};

export const logout = () => {
	return api_logout()
	.then(() => {
		store.claims = null;
	});
};
