
export const login = (name, password) => fetch("api/login", {
    method: "POST",
    body: JSON.stringify({
        name: name,
        password: password
    })
})
.then(r => r.status == 200);

export const logout = () => fetch("api/login", {
    method: "DELETE"
});

export const get_claims = () => fetch("api/login").then(r => r.status == 200 ? r.json() : null);
