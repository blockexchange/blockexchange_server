
export const find_by_keywords = keywords => m.request({
	method: "POST",
	data: {
		keywords: keywords
	},
	url: "api/searchschema"
});

export const find_recent = count => m.request({
	url: `api/searchrecent/${count}`
});

export const get_by_user_and_schemaname = (user_name, schema_name) => m.request({
	url: `/api/search/schema/byname/${user_name}/${schema_name}`
});
