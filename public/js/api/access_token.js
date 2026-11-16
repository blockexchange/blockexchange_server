
export const get_access_tokens = () => fetch(`/api/accesstoken`).then(r => r.json());

export const create_access_token = at => fetch(`/api/accesstoken`, {
    method: "POST",
    body: JSON.stringify(at)
})
.then(r => r.json());

export const delete_access_token = uid => fetch(`/api/accesstoken/${uid}`, {
    method: "DELETE"
})
.then(r => r.json());
