import { fetch_json_checked } from "./util";
import type { ErrorResponse, SchemaModDependency, SchemaModDependencyRequest } from "./types";

// all endpoints require the schema owner or an admin, rejects with an ErrorResponse otherwise

const base_url = (schema_uid: string) => `/api/schema/${encodeURIComponent(schema_uid)}/dependency`;

export const get_schema_mod_dependencies = (schema_uid: string) => fetch_json_checked<SchemaModDependency[], ErrorResponse>(base_url(schema_uid));

export const create_schema_mod_dependency = (schema_uid: string, req: SchemaModDependencyRequest) => fetch_json_checked<SchemaModDependency, ErrorResponse>(base_url(schema_uid), {
	method: "POST",
	body: JSON.stringify(req)
});

export const update_schema_mod_dependency = (schema_uid: string, uid: string, req: SchemaModDependencyRequest) => fetch_json_checked<SchemaModDependency, ErrorResponse>(`${base_url(schema_uid)}/${encodeURIComponent(uid)}`, {
	method: "PUT",
	body: JSON.stringify(req)
});

// the server answers with 204 and an empty body on success
export const delete_schema_mod_dependency = async (schema_uid: string, uid: string): Promise<void> => {
	const r = await fetch(`${base_url(schema_uid)}/${encodeURIComponent(uid)}`, {
		method: "DELETE"
	});
	if (r.status != 204) {
		throw await r.json() as ErrorResponse;
	}
};
