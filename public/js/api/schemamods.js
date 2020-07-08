
export const get_by_uid = uid => m.request({
	method: "GET",
	url: `api/schema/${uid}/mods`
});

export const create = (uid, modstats) => m.request({
	method: "POST",
	url: `api/schema/${uid}/mods`,
	body: modstats,
	headers: {
		"Authorization": localStorage.blockexchange_token
	}
});
