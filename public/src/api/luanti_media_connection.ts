import { fetch_json_checked } from "./util";
import type { LuantiMediaConnection, LuantiMediaConnectionRequest } from "./types";

export const get_luanti_media_connections = () => fetch_json_checked<LuantiMediaConnection[]>("/api/media/luanti-connection");

export const create_luanti_media_connection = (req: LuantiMediaConnectionRequest) => fetch_json_checked<LuantiMediaConnection>("/api/media/luanti-connection", {
	method: "POST",
	body: JSON.stringify(req)
});

export const update_luanti_media_connection = (id: string, req: LuantiMediaConnectionRequest) => fetch_json_checked<LuantiMediaConnection>(`/api/media/luanti-connection/${encodeURIComponent(id)}`, {
	method: "PUT",
	body: JSON.stringify(req)
});

export const delete_luanti_media_connection = (id: string) => fetch_json_checked<boolean>(`/api/media/luanti-connection/${encodeURIComponent(id)}`, {
	method: "DELETE"
});
