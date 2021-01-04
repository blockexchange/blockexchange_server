import securefetch from './securefetch.js';

export const get = schema_id => fetch(`api/schema/${schema_id}`).then(r => r.json());

export const update = schema => securefetch(`api/schema/${schema.id}`, {
	method: "PUT",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(schema)
})
.then(r => r.json());

export const remove = schema_id => securefetch(`api/schema/${schema_id}`, {
	method: "DELETE"
});
