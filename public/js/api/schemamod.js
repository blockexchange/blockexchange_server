export const get_schemamod_count = () => fetch(`${BaseURL}/api/schemamod/count`).then(r => r.json());
