import securefetch from './securefetch.js';

export const get_by_userid = user_id => fetch(`api/collection/by-userid/${user_id}`).then(r => r.json());

export const remove = collection_id => securefetch(`api/collection/${collection_id}`, {
	method: "DELETE"
});

export const create = collection => securefetch(`api/collection`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(collection)
})
.then(r => r.json());

export const update = collection => securefetch(`api/collection/${collection.id}`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(collection)
})
.then(r => r.json());
