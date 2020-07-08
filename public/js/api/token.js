
export const request_token = (name, password) => m.request({
	method: "POST",
	body: {
		name: name,
    password: password
	},
	extract: function(xhr){
		switch (xhr.status){
			case 200:
				return xhr.responseText;
			case 401:
				throw new Error("Invalid password");
			case 404:
				throw new Error("User not found");
			case 500:
				throw new Error("Server error");
		}
	},
	url: "api/token"
});
