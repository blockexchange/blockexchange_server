import { fetch_json } from "./util";

export const get_schema_star = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/star`);
export const count_schema_stars = (schema_uid: string) => fetch_json<number>(`/api/schema/${schema_uid}/star/count`);

export const star_schema = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/star`, {
	method: "PUT"
});

export const unstar_schema = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/star`, {
	method: "DELETE"
});
