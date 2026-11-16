
export const get_tags = () => fetch(`/api/tags`).then(r => r.json());
