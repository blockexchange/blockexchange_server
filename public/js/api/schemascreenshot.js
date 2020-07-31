
export const get_all = schema_id => m.request({
	method: "GET",
	url: `api/schema/${schema_id}/screenshot`
});

export const get = (schema_id, screenshot_id) => m.request({
	method: "GET",
	url: `api/schema/${schema_id}/screenshot/${screenshot_id}`
});

export const create = (schema_id, type, title, data) => m.request({
	method: "POST",
	url: `api/schema/${schema_id}/screenshot?title=${title}`,
	body: data,
	serialize: d => d,
	headers: {
		"Authorization": localStorage.blockexchange_token,
		"Content-Type": type
	}
});

export const remove = (schema_id, screenshot_id) => m.request({
	method: "DELETE",
	url: `api/schema/${schema_id}/screenshot/${screenshot_id}`,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});
