
export const get_schema_star = schema_uid => fetch(`${BaseURL}/api/schema/${schema_uid}/star`).then(r => r.json());
export const count_schema_stars = schema_uid => fetch(`${BaseURL}/api/schema/${schema_uid}/star/count`).then(r => r.json());

export const star_schema = schema_uid => fetch(`${BaseURL}/api/schema/${schema_uid}/star`, {
    method: "PUT"
})
.then(r => r.json());

export const unstar_schema = schema_uid => fetch(`${BaseURL}/api/schema/${schema_uid}/star`, {
    method: "DELETE"
})
.then(r => r.json());
