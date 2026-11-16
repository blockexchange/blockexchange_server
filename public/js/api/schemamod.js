export const get_schemamod_count = () => fetch(`/api/schemamod/count`).then(r => r.json());
