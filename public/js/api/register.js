
export const register = rr => fetch(`${BaseURL}/api/register`, {
    method: "POST",
    body: JSON.stringify(rr)
})
.then(r => r.json());
