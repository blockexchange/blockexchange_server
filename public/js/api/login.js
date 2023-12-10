
export const login = (name, password) => fetch(`${BaseURL}/api/login`, {
    method: "POST",
    body: JSON.stringify({
        name: name,
        password: password
    })
})
.then(r => r.status == 200);

export const logout = () => fetch(`${BaseURL}/api/login`, {
    method: "DELETE"
});

export const get_claims = () => fetch(`${BaseURL}/api/login`).then(r => r.status == 200 ? r.json() : null);
