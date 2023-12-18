
export const get_schema_star = schema_id => fetch(`${BaseURL}/api/schema/${schema_id}/star`).then(r => r.json());
export const count_schema_stars = schema_id => fetch(`${BaseURL}/api/schema/${schema_id}/star/count`).then(r => r.json());

export const star_schema = schema_id => fetch(`${BaseURL}/api/schema/${schema_id}/star`, {
    method: "PUT"
})
.then(r => r.json());

export const unstar_schema = schema_id => fetch(`${BaseURL}/api/schema/${schema_id}/star`, {
    method: "DELETE"
})
.then(r => r.json());
