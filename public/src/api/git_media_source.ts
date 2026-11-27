import { fetch_json_checked } from "./util";
import type { GitMediaSource, GitMediaSourceRequest } from "./types";

export const get_git_media_sources = () => fetch_json_checked<GitMediaSource[]>("/api/media/git-source");

export const create_git_media_source = (req: GitMediaSourceRequest) => fetch_json_checked<GitMediaSource>("/api/media/git-source", {
	method: "POST",
	body: JSON.stringify(req)
});

export const update_git_media_source = (id: string, req: GitMediaSourceRequest) => fetch_json_checked<GitMediaSource>(`/api/media/git-source/${encodeURIComponent(id)}`, {
	method: "PUT",
	body: JSON.stringify(req)
});

export const delete_git_media_source = (id: string) => fetch_json_checked<boolean>(`/api/media/git-source/${encodeURIComponent(id)}`, {
	method: "DELETE"
});

// imports the media files right away, resolves with the updated status
export const run_git_media_source = (id: string) => fetch_json_checked<GitMediaSource>(`/api/media/git-source/${encodeURIComponent(id)}/run`, {
	method: "POST"
});
