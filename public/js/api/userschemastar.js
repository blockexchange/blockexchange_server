export const remove = schema_id => m.request({
	method: "DELETE",
	url: `api/schema/${schema_id}/star`,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});

export const add = schema_id => m.request({
	method: "PUT",
	url: `api/schema/${schema_id}/star`,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});

export const get_all = schema_id => m.request({
	url: `api/schema/${schema_id}/star`,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});
