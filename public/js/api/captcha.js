
export const create_captcha = () => fetch(`api/captcha`).then(r => r.text());