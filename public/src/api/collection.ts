import { fetch_json } from "./util";
import type { Collection } from "./types";

export const get_collection = (collection_uid: string) => fetch_json<Collection>(`/api/collection/${collection_uid}`);

export const get_collections_by_username = (username: string) => fetch_json<Collection[]>(`/api/collection/by-username/${username}`);

export const create_collection = (c: Collection) => fetch_json<Collection>(`/api/collection`, {
	method: "POST",
	body: JSON.stringify(c)
});

export const update_collection = (c: Collection) => fetch_json<Collection>(`/api/collection`, {
	method: "PUT",
	body: JSON.stringify(c)
});

export const delete_collection = (collection_uid: string) => fetch_json<boolean>(`/api/collection/${collection_uid}`, {
	method: "DELETE"
});
