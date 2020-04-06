
export const get_by_uid = uid => m.request({
	method: "GET",
	url: `api/schema/${uid}/mods`
});
