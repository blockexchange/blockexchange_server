import securefetch from './securefetch.js';

export const get = schema_id => fetch(`api/schematag/${schema_id}`).then(r => r.json());

export const add = (schema_id, tag_id) => securefetch(`api/schematag/${schema_id}/${tag_id}`, {
	method: "PUT"
});

export const remove = (schema_id, tag_id) => securefetch(`api/schematag/${schema_id}/${tag_id}`, {
	method: "DELETE"
});
