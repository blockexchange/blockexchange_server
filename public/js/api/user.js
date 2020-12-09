
export const get_all = () => fetch("api/user").then(r => r.json());
