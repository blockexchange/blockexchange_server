
export const request_token = (name, password) => fetch("api/token", {
    method: "POST",
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ name: name, password: password })
  })
  .then(r => r.text());
