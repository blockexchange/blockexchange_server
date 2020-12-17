const bcrypt = require('bcryptjs');
const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const createjwt = require("../util/createjwt");
const logger = require("../logger");
const app = require("../app");
const user_dao = require("../dao/user");
const tokenmiddleware = require("../middleware/token");
const tokencheck = tokenmiddleware();

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

	if (req.body.password === ""){
		return res.status(403).json({
			message: "Empty password not allowed"
		});
	}

	const success = bcrypt.compareSync(req.body.password, user.hash);
	if (!success) {
		return res.status(401).json({
			message: "Invalid password"
		});
	}

	const token = createjwt(user);
	res.send(token);
});

// returns a new token with the updated user data
app.post("/api/token/refresh", tokencheck, async function(req, res){
	const user = await user_dao.get_by_id(req.claims.user_id);
	const token = createjwt(user);
	res.send(token);
});
