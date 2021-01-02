import securefetch from './securefetch.js';

export const get_all = () => securefetch("api/access_token").then(r => r.json());

export const create = (name, expires) => securefetch(`api/access_token`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify({
		name: name,
		expires: expires
	})
})
.then(r => r.json());

export const remove = id => securefetch(`api/access_token/${id}`, {
	method: "DELETE"
});
