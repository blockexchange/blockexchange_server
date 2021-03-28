import securefetch from './securefetch.js';

export const get_all = () => fetch("api/tag").then(r => r.json());

export const remove = tag_id => securefetch(`api/tag/${tag_id}`, {
	method: "DELETE"
});

export const create = tag => securefetch(`api/tag`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(tag)
})
.then(r => r.json());

export const update = tag => securefetch(`api/tag/${tag.id}`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(tag)
})
.then(r => r.json());
