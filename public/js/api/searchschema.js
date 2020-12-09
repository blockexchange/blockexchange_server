
export const find_recent = count => fetch(`api/searchrecent/${count}`).then(r => r.json());

export const search = params => fetch("api/searchschema", {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(params)
})
.then(r => r.json());
