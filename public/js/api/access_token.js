
export const get_access_tokens = () => fetch(`${BaseURL}/api/accesstoken`).then(r => r.json());

export const create_access_token = at => fetch(`${BaseURL}/api/accesstoken`, {
    method: "POST",
    body: JSON.stringify(at)
})
.then(r => r.json());

export const delete_access_token = id => fetch(`${BaseURL}/api/accesstoken/${id}`, {
    method: "DELETE"
})
.then(r => r.json());
