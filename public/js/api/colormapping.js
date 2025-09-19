
export const get_colormapping = () => fetch(`${BaseURL}/api/colormapping`).then(r => r.json());
