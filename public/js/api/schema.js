
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

export const schema_update = schema => fetch(`${BaseURL}/api/schema/${schema.uid}`, {
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

export const schema_set_tags = (schema_uid, tags) => fetch(`${BaseURL}/api/schema/${schema_uid}/tags`, {
    method: "POST",
    body: JSON.stringify(tags)
})
.then(r => r.json());

export const schema_update_screenshot = schema_uid => fetch(`/api/schema/${schema_uid}/screenshot/update`, {
    method: "POST"
})
.then(r => r.json());

export const schema_update_info = schema_uid => fetch(`/api/schema/${schema_uid}/update`, {
    method: "POST"
})
.then(r => r.json());


export const schema_delete = schema_uid => fetch(`${BaseURL}/api/schema/${schema_uid}`, {
    method: "DELETE",
})
.then(r => r.json());

export const schema_update_mods = schema_uid => fetch(`/api/schema/${schema_uid}/mods/update`, {
    method: "POST"
})
.then(r => r.json());
