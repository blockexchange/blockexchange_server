
export const get_by_uid = uid => m.request({
	method: "GET",
	url: `api/schema/${uid}`
});

export const update = schema => m.request({
	method: "PUT",
	data: schema,
	url: `api/schema/${schema.id}`
});
