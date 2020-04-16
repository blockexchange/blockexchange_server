
export const register = (name, password, mail) => m.request({
	method: "POST",
	data: {
		name: name,
    password: password,
    mail: mail
	},
	url: "api/register"
});
