import securefetch from './securefetch.js';
import memoize from '../util/memoize.js';

export const get_all = memoize(() => fetch("api/user").then(r => r.json()));

export const update = user => securefetch(`api/user/${user.id}`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(user)
})
.then(r => r.json());

export const validate_username = username => securefetch(`api/validate_username`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify({ username: username })
})
.then(r => r.json());
