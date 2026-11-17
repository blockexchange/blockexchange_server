// fetches the url and parses the response body as json
export const fetch_json = <T>(url: string, init?: RequestInit): Promise<T> => fetch(url, init).then(r => r.json() as Promise<T>);

// same as fetch_json but rejects with the parsed json body on a non-200 status
export const fetch_json_checked = <T, E = unknown>(url: string, init?: RequestInit): Promise<T> => fetch(url, init)
.then(r => {
	if (r.status == 200) {
		return r.json() as Promise<T>;
	} else {
		return r.json().then((msg: E) => Promise.reject(msg));
	}
});
