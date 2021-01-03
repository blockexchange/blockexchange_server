import securefetch from './securefetch.js';

export const get_by_collectionid = collection_id => fetch(`api/collection_schema/by-collectionid/${collection_id}`).then(r => r.json());

export const remove = (collection_id, schema_id) => securefetch(`api/collection_schema/${collection_id}/${schema_id}`, {
	method: "DELETE"
});

export const create = collection_schema => securefetch(`api/collection_schema`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(collection_schema)
})
.then(r => r.json());
