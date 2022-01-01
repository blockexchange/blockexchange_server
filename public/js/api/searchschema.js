
export const find_recent = (limit, offset) => fetch(`api/searchrecent?limit=${limit}&offset=${offset || 0}`).then(r => r.json());

export const search = (params, limit, offset) => fetch(`api/searchschema?limit=${limit}&offset=${offset || 0}`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(params)
})
.then(r => r.json());

export const search_by_user_and_schemaname = (username, schemaname) => fetch(`api/search/schema/byname/${username}/${schemaname}`)
	.then(r => r.json());
