
export const get_schemapart_by_ofset = (uid,x,y,z) => fetch(`${BaseURL}/api/schemapart/${uid}/${x}/${y}/${z}`).then(r => r.json());
