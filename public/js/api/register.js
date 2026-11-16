
export const register = rr => fetch(`/api/register`, {
    method: "POST",
    body: JSON.stringify(rr)
})
.then(r => r.json());
