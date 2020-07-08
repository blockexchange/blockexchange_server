
export const get = (schema_id, block_x, block_y, block_z) => m.request({
	method: "GET",
	url: `api/schemapart/${schema_id}/${block_x}/${block_y}/${block_z}`
});

export const create = part => m.request({
	method: "POST",
  url: `api/schemapart`,
	body: part,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});
