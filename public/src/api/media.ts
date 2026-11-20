import { fetch_json, fetch_json_checked } from "./util";
import type { Mediafile, MediafileInfo, Nodedefinition } from "./types";

export interface MediaQuery {
	q?: string;
	limit?: number;
	offset?: number;
}

const to_params = (query: MediaQuery) => new URLSearchParams(
	Object.entries(query)
		.filter(([, v]) => v !== undefined && v !== "")
		.map(([k, v]) => [k, String(v)])
).toString();

export const get_nodedefinitions = (query: MediaQuery = {}) => fetch_json<Nodedefinition[]>(`/api/media/nodedef?${to_params(query)}`);

// returns the definitions of the given node names (max 1000), unknown names are skipped
export const get_nodedefinitions_batch = (names: string[]) => fetch_json_checked<Nodedefinition[]>(`/api/media/nodedef-batch`, {
	method: "POST",
	body: JSON.stringify(names)
});

export const count_nodedefinitions = (q?: string) => fetch_json<number>(`/api/media/nodedef-count?${to_params({ q })}`);

export const delete_nodedefinition = (name: string) => fetch_json_checked<boolean>(`/api/media/nodedef/${encodeURIComponent(name)}`, {
	method: "DELETE"
});

export const get_mediafiles = (query: MediaQuery = {}) => fetch_json<MediafileInfo[]>(`/api/media/mediafile?${to_params(query)}`);

// returns the media files with their data for the given file names (max 1000), unknown names are skipped
export const get_mediafiles_batch = (names: string[]) => fetch_json_checked<Mediafile[]>(`/api/media/mediafile-batch`, {
	method: "POST",
	body: JSON.stringify(names)
});

export const count_mediafiles = (q?: string) => fetch_json<number>(`/api/media/mediafile-count?${to_params({ q })}`);

export const get_mediafile_raw_url = (name: string) => `/api/media/mediafile/${encodeURIComponent(name)}?raw=true`;

export const delete_mediafile = (name: string) => fetch_json_checked<boolean>(`/api/media/mediafile/${encodeURIComponent(name)}`, {
	method: "DELETE"
});
