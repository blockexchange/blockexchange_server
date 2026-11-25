import { fetch_json } from "./util";
import type { SchemaSearchRequest, SchemaSearchResponse } from "./types";

export const get_schema_star = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/star`);
export const count_schema_stars = (schema_uid: string) => fetch_json<number>(`/api/schema/${schema_uid}/star/count`);

export const star_schema = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/star`, {
	method: "PUT"
});

export const unstar_schema = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/star`, {
	method: "DELETE"
});

// schematics starred by the logged-in user (stars are private, the user is taken from the token)
export const starred_schema_search = (search: SchemaSearchRequest) => fetch_json<SchemaSearchResponse[]>(`/api/search/schema/starred`, {
	method: "POST",
	body: JSON.stringify(search)
});

export const starred_schema_count = (search: SchemaSearchRequest) => fetch_json<number>(`/api/count/schema/starred`, {
	method: "POST",
	body: JSON.stringify(search)
});
