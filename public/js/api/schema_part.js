
function responseHandler(r) {
    if (r.status == 200) {
        return r.json();
    }
    return Promise.resolve(null);
}

export const get_schemapart_by_offset = (uid,x,y,z) => fetch(`${BaseURL}/api/schemapart/${uid}/${x}/${y}/${z}`).then(responseHandler);

export const get_next_schemapart_by_offset = (uid,x,y,z) => fetch(`${BaseURL}/api/schemapart_next/${uid}/${x}/${y}/${z}`).then(responseHandler);

export const count_schemaparts_by_mtime = (uid, mtime) => fetch(`${BaseURL}/api/schemapart_count/by-mtime/${uid}/${mtime}`).then(responseHandler);

export const get_next_schemapart_by_mtime = (uid, mtime) => fetch(`${BaseURL}/api/schemapart_next/by-mtime/${uid}/${mtime}`).then(responseHandler);