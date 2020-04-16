
export const get_token = (name, password) => m.request({
	method: "POST",
	data: {
		name: name,
    password: password
	},
	url: "api/token"
});
