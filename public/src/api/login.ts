import type { Claims } from "./types";

export const login = async (name: string, password: string) => {
	const r = await fetch(`/api/login`, {
		method: "POST",
		body: JSON.stringify({
			name: name,
			password: password
		})
	});
	return r.status == 200;
};

export const logout = () => fetch(`/api/login`, {
	method: "DELETE"
});

export const get_claims = async (renew?: boolean): Promise<Claims | null> => {
	const r = await fetch(`/api/login${renew ? "?renew=true" : ""}`);
	return r.status == 200 ? await r.json() as Claims : null;
};
