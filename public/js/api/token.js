
export const request_token = (name, password) => fetch("api/token", {
    method: "POST",
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ name: name, password: password })
  })
  .then(r => {
    switch (r.status){
      case 404:
        throw new Error("User not found");
      case 200:
        return r.text();
      default:
        throw new Error("unknown error");
    }
  });
