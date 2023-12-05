
export const get_login = () => fetch("api/login").then(r => r.json());

export const do_login = (name, password) => fetch("api/login", {
    method: "POST",
    body: JSON.stringify({
        name: name,
        password: password
    })
}).then(r => r.json());
