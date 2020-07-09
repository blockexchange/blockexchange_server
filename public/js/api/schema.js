
export const get_by_id = id => m.request({
	method: "GET",
	url: `api/schema/${id}`
});

export const update = schema => m.request({
	method: "PUT",
	url: `api/schema/${schema.id}`,
	body: schema,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});

export const create = schema => m.request({
	method: "POST",
	url: `api/schema`,
	body: schema,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});

export const complete = schema => m.request({
	method: "POST",
	url: `api/schema/${schema.id}/complete`,
	body: schema,
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
