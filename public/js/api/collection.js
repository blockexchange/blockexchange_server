
export const get_collection = collection_uid => fetch(`${BaseURL}/api/collection/${collection_uid}`).then(r => r.json());

export const get_collections_by_username = username => fetch(`${BaseURL}/api/collection/by-username/${username}`).then(r => r.json());

export const create_collection = c => fetch(`${BaseURL}/api/collection`, {
    method: "POST",
    body: JSON.stringify(c)
})
.then(r => r.json());

export const update_collection = c => fetch(`${BaseURL}/api/collection`, {
    method: "PUT",
    body: JSON.stringify(c)
})
.then(r => r.json());

export const delete_collection = collection_uid => fetch(`${BaseURL}/api/collection/${collection_uid}`, {
    method: "DELETE"
})
.then(r => r.json());
