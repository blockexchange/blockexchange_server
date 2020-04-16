
export const request_token = (name, password) => m.request({
	method: "POST",
	data: {
		name: name,
    password: password
	},
	url: "api/token"
});
