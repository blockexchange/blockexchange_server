
export const register = (name, password) => fetch(`api/register`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify({
		name: name,
		password: password
	})
})
.then(r => r.json());
