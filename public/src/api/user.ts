import { fetch_json } from "./util";
import type { ChangePassword, User, UserSearch } from "./types";

export const count_users = () => fetch_json<number>(`/api/user-count`);

export const search_users = (search: UserSearch) => fetch_json<User[]>(`/api/user-search`, {
	method: "POST",
	body: JSON.stringify(search)
});

export const save_user = (user: User) => fetch_json<User>(`/api/user/${user.uid}`, {
	method: "POST",
	body: JSON.stringify(user)
});

export const get_user = (id: string) => fetch_json<User>(`/api/user/${id}`);

export const count_user_schema_stars = (uid: string) => fetch_json<number>(`/api/user/${uid}/schemastars`);

export const change_password = async (userid: string, chpwd: ChangePassword) => {
	const r = await fetch(`/api/user/${userid}/changepassword`, {
		method: "POST",
		body: JSON.stringify(chpwd)
	});
	return r.status == 200;
};

export const unlink_oauth = async (userid: string, chpwd: ChangePassword) => {
	const r = await fetch(`/api/user/${userid}/unlink-oauth`, {
		method: "POST",
		body: JSON.stringify(chpwd)
	});
	return r.status == 200;
};
