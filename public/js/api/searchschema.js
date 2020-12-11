
export const find_recent = count => fetch(`api/searchrecent/${count}`).then(r => r.json());

export const search = params => fetch("api/searchschema", {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(params)
})
.then(r => r.json());

export const search_by_user_and_schemaname = (username, schemaname) => fetch(`api/search/schema/byname/${username}/${schemaname}`)
	.then(r => r.json());
