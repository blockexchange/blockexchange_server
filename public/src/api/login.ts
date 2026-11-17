import type { Claims } from "./types";

export const login = (name: string, password: string) => fetch(`/api/login`, {
	method: "POST",
	body: JSON.stringify({
		name: name,
		password: password
	})
})
.then(r => r.status == 200);

export const logout = () => fetch(`/api/login`, {
	method: "DELETE"
});

export const get_claims = (renew?: boolean): Promise<Claims | null> => fetch(`/api/login${renew ? "?renew=true" : ""}`)
.then(r => r.status == 200 ? r.json() as Promise<Claims> : null);
