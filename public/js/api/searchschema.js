
export const search = (params, limit, offset) => fetch(`api/searchschema?limit=${limit}&offset=${offset || 0}`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(params)
})
.then(r => r.json());
