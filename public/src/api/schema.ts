import { fetch_json, fetch_json_checked } from "./util";
import type { Schema, SchemaSearchRequest, SchemaSearchResponse, SchemaUpdateError } from "./types";

export const schema_search = (search: SchemaSearchRequest) => fetch_json<SchemaSearchResponse[]>(`/api/search/schema`, {
	method: "POST",
	body: JSON.stringify(search)
});

export const schema_count = (search: SchemaSearchRequest) => fetch_json<number>(`/api/count/schema`, {
	method: "POST",
	body: JSON.stringify(search)
});

export const schema_update = (schema: Schema) => fetch_json_checked<Schema, SchemaUpdateError>(`/api/schema/${schema.uid}`, {
	method: "PUT",
	body: JSON.stringify(schema)
});

export const schema_set_tags = (schema_uid: string, tags: string[]) => fetch_json<boolean>(`/api/schema/${schema_uid}/tags`, {
	method: "POST",
	body: JSON.stringify(tags)
});

export const schema_update_screenshot = (schema_uid: string) => fetch_json<unknown>(`/api/schema/${schema_uid}/screenshot/update`, {
	method: "POST"
});

export const schema_update_info = (schema_uid: string) => fetch_json<Schema>(`/api/schema/${schema_uid}/update`, {
	method: "POST"
});

export const schema_delete = (schema_uid: string) => fetch_json<boolean>(`/api/schema/${schema_uid}`, {
	method: "DELETE",
});

export const schema_update_mods = (schema_uid: string) => fetch_json<string[]>(`/api/schema/${schema_uid}/mods/update`, {
	method: "POST"
});
