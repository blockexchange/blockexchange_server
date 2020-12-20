import securefetch from './securefetch.js';

export const request_token = (name, password) => fetch("api/token", {
    method: "POST",
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({ name: name, password: password })
  })
  .then(r => {
		if (r.status != 200){
			return r.json()
			.then(response => {
				throw new Error(response.message);
			});
		} else {
			return r.text();
		}
  });

export const create = options => securefetch("api/token/create", {
		method: "POST",
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(options)
	})
	.then(r => r.text());
