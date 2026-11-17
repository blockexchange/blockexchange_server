import { fetch_json } from "./util";
import type { AccessToken } from "./types";

export const get_access_tokens = () => fetch_json<AccessToken[]>(`/api/accesstoken`);

export const create_access_token = (at: Pick<AccessToken, "name" | "expires">) => fetch_json<AccessToken>(`/api/accesstoken`, {
	method: "POST",
	body: JSON.stringify(at)
});

export const delete_access_token = (uid: string) => fetch_json<boolean>(`/api/accesstoken/${uid}`, {
	method: "DELETE"
});
