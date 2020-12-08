const axios = require('axios');

const app = require("../app");
const logger = require("../logger");
const user_dao = require("../dao/user");
const createjwt = require("../util/createjwt");

app.get('/api/oauth_callback', function(req, res){
  logger.debug("GET /api/oauth_callback", req.query);

  const data = {
    client_id: process.env.GITHUB_APP_ID,
    client_secret: process.env.GITHUB_APP_SECRET,
    code: req.query.code
  };

  const options = {
    headers: {
      "Accept": "application/json"
    }
  };

	let user_info;

  axios.post("https://github.com/login/oauth/access_token", data, options)
  .then(r => {
    if (!r.data.access_token){
      console.log(r);
      throw new Error("No access_token received!");
    }
    return axios.get("https://api.github.com/user", {
      headers: {
        "Authorization": "Bearer " + r.data.access_token
      }
    });
  })
  .then(r => {
    user_info = r.data;
    console.log(user_info);
		return user_dao.get_by_name(user_info.login);
    // user_info.login / avatar_url / name / email
	})
	.then(user => {
		if (user){
			if (user.type != "GITHUB") {
				throw new Error("User already exists with another type!");
			}
			// valid github user
			return user;
		} else {
			// create new user
			return user_dao.create({
				name: user_info.login,
				role: "MEMBER",
				type: "GITHUB",
				hash: "", // no local password
				mail: user_info.email
			});
		}
	})
	.then(user => {
		const token = createjwt(user);
		res.redirect(process.env.BASE_URL + "/#/oauth?token=" + token);
	})
  .catch(e => {
    console.log(e);
    res.send(e.message);
  });
});
