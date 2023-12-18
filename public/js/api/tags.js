
export const get_tags = () => fetch(`${BaseURL}/api/tags`).then(r => r.json());
