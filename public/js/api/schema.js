import securefetch from './securefetch.js';

export const get = (username, schemaname) => fetch(`api/schema/${username}/${schemaname}`).then(r => r.json());

export const update = schema => securefetch(`api/schema`, {
	method: "PUT",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(schema)
})
.then(r => r.json());

export const updateInfo = (username, schemaname) => securefetch(`api/schema/${username}/${schemaname}/update`, {
	method: "POST"
})
.then(r => r.text());

export const remove = (username, schemaname) => securefetch(`api/schema/${username}/${schemaname}`, {
	method: "DELETE"
});
