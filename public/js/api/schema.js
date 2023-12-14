
export const schema_search = search => fetch(`${BaseURL}/api/search/schema`, {
    method: "POST",
    body: JSON.stringify(search)
})
.then(r => r.json());
