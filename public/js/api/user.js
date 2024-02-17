export const count_users = () => fetch(`${BaseURL}/api/user-count`).then(r => r.json());

export const search_users = search => fetch(`${BaseURL}/api/user-search`, {
    method: "POST",
    body: JSON.stringify(search)
})
.then(r => r.json());

export const save_user = user => fetch(`${BaseURL}/api/user/${user.uid}`, {
    method: "POST",
    body: JSON.stringify(user)
})
.then(r => r.json());

export const get_user = id => fetch(`${BaseURL}/api/user/${id}`).then(r => r.json());

export const count_user_schema_stars = uid => fetch(`${BaseURL}/api/user/${uid}/schemastars`).then(r => r.json());
