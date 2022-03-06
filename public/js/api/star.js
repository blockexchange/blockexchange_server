import securefetch from './securefetch.js';

export const get = (schema_id, user_id) => fetch(`api/schema/${schema_id}/star${user_id ? '?user_id='+user_id : ''}`).then(r => r.json());

export const remove = schema_id => securefetch(`api/schema/${schema_id}/star`, {
	method: "DELETE"
});

export const create = schema_id => securefetch(`api/schema/${schema_id}/star`, {
	method: "POST"
});
