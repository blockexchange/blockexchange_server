import securefetch from './securefetch.js';

export const count = schema_id => fetch(`api/schema/${schema_id}/star`)
.then(r => r.text())
.then(count => +count);

export const remove = schema_id => securefetch(`api/schema/${schema_id}/star`, {
	method: "DELETE"
});

export const create = schema_id => securefetch(`api/schema/${schema_id}/star`, {
	method: "POST"
});
