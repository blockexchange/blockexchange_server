
const cache = {};

export const get_by_schemaid = schema_id => {
	if (cache[schema_id]){
		return cache[schema_id];
	}
	fetch(`api/schema/${schema_id}/screenshot`)
	.then(r => {
		const json = r.json();
		cache[schema_id] = json;
		return json;
	});
};
