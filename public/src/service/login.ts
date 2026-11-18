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

export const check_login = async (renew?: boolean) => {
	const c = await fetch_claims(renew);
	store.claims = c;
	return c;
};

export const login = async (name: string, password: string) => {
	const success = await api_login(name, password);
	return success ? await check_login() : null;
};

export const logout = async () => {
	await api_logout();
	store.claims = null;
};
