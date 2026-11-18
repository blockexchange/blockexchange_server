export const create_captcha = async () => {
	const r = await fetch(`/api/captcha`);
	return await r.text();
};
