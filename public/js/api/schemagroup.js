
export const get_by_id = id => m.request({
	method: "GET",
	url: `api/schemagroup/${id}`
});
