
export const get_all = () => fetch("api/user").then(r => r.json());

export const update = user => fetch(`api/user/${user.id}`, {
	method: "POST",
	headers: {
		'Content-Type': 'application/json'
	},
	body: JSON.stringify(user)
})
.then(r => r.json());
