
export const schema_search = search => fetch(`${BaseURL}/api/search/schema`, {
    method: "POST",
    body: JSON.stringify(search)
})
.then(r => r.json());

export const schema_count = search => fetch(`${BaseURL}/api/count/schema`, {
    method: "POST",
    body: JSON.stringify(search)
})
.then(r => r.json());

export const get_schema_by_name = (username, name) => fetch(`${BaseURL}/api/search/schema/byname/${username}/${name}`).then(r => r.json());

export const schema_update = schema => fetch(`${BaseURL}/api/schema/${schema.id}`, {
    method: "PUT",
    body: JSON.stringify(schema)
})
.then(r => {
    if (r.status == 200) {
        return r.json();
    } else {
        return r.json().then(msg => Promise.reject(msg));
    }
});