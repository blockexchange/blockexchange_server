// fetches the url and parses the response body as json
export const fetch_json = async <T>(url: string, init?: RequestInit): Promise<T> => {
	const r = await fetch(url, init);
	return await r.json() as T;
};

// same as fetch_json but rejects with the parsed json body on a non-200 status
export const fetch_json_checked = async <T, E = unknown>(url: string, init?: RequestInit): Promise<T> => {
	const r = await fetch(url, init);
	if (r.status != 200) {
		throw await r.json() as E;
	}
	return await r.json() as T;
};
