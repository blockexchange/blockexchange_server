
import { get_claims as fetch_claims, login as api_login, logout as api_logout } from '../api/login.js';

const store = Vue.reactive({
    claims: null
});

export const is_logged_in = () => store.claims != null;
export const get_claims = () => store.claims;

export const check_login = () => fetch_claims().then(c => {
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