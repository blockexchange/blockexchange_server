import { fetch_json, fetch_json_checked } from "./util";
import type { Tag, TagRequest } from "./types";

export const get_tags = () => fetch_json<Tag[]>(`/api/tags`);

export const create_tag = (req: TagRequest) => fetch_json_checked<Tag>(`/api/tags`, {
	method: "POST",
	body: JSON.stringify(req)
});

export const update_tag = (uid: string, req: TagRequest) => fetch_json_checked<Tag>(`/api/tags/${encodeURIComponent(uid)}`, {
	method: "PUT",
	body: JSON.stringify(req)
});
