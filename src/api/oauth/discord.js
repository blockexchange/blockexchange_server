const axios = require('axios');
const qs = require('querystring');

const app = require("../../app");
const logger = require("../../logger");
const user_dao = require("../../dao/user");
const setupuser = require("../../util/setupuser");

const { MANAGEMENT, UPLOAD, OVERWRITE } = require("../../permissions");
const createjwt = require("../../util/createjwt");

app.get('/api/oauth_callback/discord', function(req, res){
  logger.debug("GET /api/oauth_callback/discord", req.query);

  const data = {
    client_id: process.env.DISCORD_APP_ID,
    client_secret: process.env.DISCORD_APP_SECRET,
    code: req.query.code,
		grant_type: 'authorization_code',
		redirect_uri: process.env.BASE_URL + "/api/oauth_callback/discord",
    'scope': 'identify email connections'
  };

  const options = {
    headers: {
      "Accept": "application/json",
			"Content-Type": "application/x-www-form-urlencoded"
    }
  };

	let user_info;
	axios.post("https://discord.com/api/oauth2/token", qs.stringify(data), options)
  .then(r => {
		console.log(r.data);
    if (!r.data.access_token){
      console.log(r);
      throw new Error("No access_token received!");
    }

		return axios.get("https://discord.com/api/users/@me", {
      headers: {
        "Authorization": "Bearer " + r.data.access_token
      }
    })
		.then(r => {
			user_info = r.data;
			console.log(user_info);
			// id, username, email, verified
			return user_dao.get_by_external_id(user_info.id);
		})
		.then(user => {
			if (user){
				if (user.type != "DISCORD") {
					throw new Error("User already exists with another type!");
				}
				// valid user
				return user;
			} else {
				// create new user
				return user_dao.create({
					name: user_info.username,
					type: "DISCORD",
					hash: "", // no local password
					mail: user_info.email,
					external_id: user_info.id
				})
				.then(setupuser);
			}
		});
  })
	.then(user => {
		const token = createjwt(user, [UPLOAD, OVERWRITE, MANAGEMENT]);
		res.redirect(process.env.BASE_URL + "/#/oauth/" + token);
	})
	.catch(e => {
		console.log(e);
		res.send(e.message);
	});

});
