
export const import_schematic = (buf, filename) => fetch(`/api/import/${filename}`, {
    method: "POST",
    body: buf
})
.then(r => {
    if (r.status == 200) {
        return r.json();
    } else {
        return r.json().then(msg => Promise.reject(msg));
    }
});