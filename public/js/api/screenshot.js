export const get_by_schemaid = schema_id => fetch(`api/schema/${schema_id}/screenshot`)
	.then(r => r.json());
