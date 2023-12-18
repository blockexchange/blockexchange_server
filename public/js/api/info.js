
export const get_info = () => fetch(`${BaseURL}/api/info`).then(r => r.json());
