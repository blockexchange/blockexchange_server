
export const create_captcha = () => fetch(`${BaseURL}/api/captcha`).then(r => r.text());