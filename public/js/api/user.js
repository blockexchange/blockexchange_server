import securefetch from './securefetch.js';

export const get = (limit, offset) => fetch(`api/user?limit=${limit}&offset=${offset || 0}`).then(r => r.json());

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
	body: JSON.stringify({ name: username })
})
.then(r => r.json());
