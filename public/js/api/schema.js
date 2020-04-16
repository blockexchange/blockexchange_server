
export const get_by_uid = uid => m.request({
	method: "GET",
	url: `api/schema/${uid}`
});

export const update = schema => m.request({
	method: "PUT",
	url: `api/schema/${schema.id}`,
	data: schema,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});

export const remove = schema => m.request({
	method: "DELETE",
	url: `api/schema/${schema.id}`,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});
