import store from '../store/login.js';

export default function(url, options) {
	options.headers = options.headers || {};
	options.headers.Authorization = store.token;

	return fetch(url, options);
}
