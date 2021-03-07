
export const get_tags = () => fetch("api/tag").then(r => r.json());
