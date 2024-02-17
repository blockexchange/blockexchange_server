
import { get_claims as fetch_claims, login as api_login, logout as api_logout } from '../api/login.js';

const store = Vue.reactive({
    claims: null
});

export const is_logged_in = () => store.claims != null;
export const get_claims = () => store.claims;

export const get_user_uid = () => store.claims ? store.claims.user_uid : null;
export const get_username = () => store.claims ? store.claims.username : null;
export const has_permission = permission => store.claims ? store.claims.permissions.includes(permission) : false;

export const check_login = renew => fetch_claims(renew).then(c => {
    store.claims = c;
    return c;
});

export const login = (name, password) => {
    return api_login(name, password)
    .then(success => {
        if (success) {
            return check_login();
        }
    });
};

export const logout = () => {
    return api_logout()
    .then(() => {
        store.claims = null;
    });
};