const bcrypt = require('bcryptjs');
const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const createjwt = require("../util/createjwt");
const logger = require("../logger");
const app = require("../app");

const user_dao = require("../dao/user");
const access_token_dao = require("../dao/access_token");

// data='{"name": "xyz", "password": "abc"}'
// curl -X POST 127.0.0.1:8080/api/token --data "${data}" -H "Content-Type: application/json"
app.post('/api/token', jsonParser, async function(req, res){
	logger.debug("POST /api/token", req.body.name);

	const user = await user_dao.get_by_name(req.body.name);
	if (!user) {
		return res.status(404).json({
			message: "User not found"
		});
	}

	if (req.body.password){
		// login with username/password (from web-ui)
		const success = bcrypt.compareSync(req.body.password, user.hash);
		if (!success) {
			return res.status(401).json({
				message: "Invalid password"
			});
		}
		const token = createjwt(user, {
			audience: "management"
		});
		res.send(token);

	} else if (req.body.access_token){
		// login with access_token (from ingame)
		const access_token = await access_token_dao.find_by_username_and_token(req.body.name, req.body.access_token);

		if (!access_token){
			return res.status(401).json({
				message: "token not found"
			});
		}

		if (access_token.expires < Date.now()){
			return res.status(403).json({
				message: "token expired"
			});
		}

		// valid token and expiration time
		const token = createjwt(user, {
			audience: "minetest",
			expiresIn: parseInt((access_token.expires - Date.now()) / 1000)
		});
		res.send(token);

	} else {
		// invalid request
		return res.status(403).json({
			message: "Empty password/access_token not allowed"
		});
	}

});
