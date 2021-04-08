import securefetch from './securefetch.js';

export const get = schema_id => fetch(`api/schema/${schema_id}/tag`).then(r => r.json());

export const add = (schema_id, tag_id) => securefetch(`api/schema/${schema_id}/tag/${tag_id}`, {
	method: "PUT"
});

export const remove = (schema_id, tag_id) => securefetch(`api/schema/${schema_id}/tag/${tag_id}`, {
	method: "DELETE"
});
