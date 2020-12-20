const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const user_dao = require("../dao/user");
const logger = require("../logger");
const tokenmiddleware = require("../middleware/token");
const tokencheck = tokenmiddleware();

// get all users
app.get('/api/user', async function(req, res){
	logger.debug("GET /api/user");

	const users = await user_dao.get_all();
	const list = users.map(user => {
		return {
			name: user.name,
			id: user.id,
			type: user.type,
			role: user.role,
			created: user.created
		};
	});

	res.json(list);
});

// validates the new username, checks for duplicates and invalid characters
app.post("/api/validate_username", jsonParser, tokencheck, async function(req, res){
	if (! /^[a-z0-9\_\-]+$/i.test(req.body.username)) {
		return res.json({
			valid: false,
			message: "Username can only contain characters,numbers,dashes and underlines"
		});
	}

	const existing_user = await user_dao.get_by_name(req.body.username);

	if (existing_user){
		return res.json({
			valid: false,
			message: "User already exists"
		});
	}

	// all checks ok
	res.json({
		valid: true
	});
});

// modify user (only name change allowed)
app.post("/api/user/:user_id", jsonParser, tokencheck, async function(req, res){
	logger.debug("POST /api/user", req.params.user_id, req.body);

	if (req.claims.role !== "ADMIN" && req.claims.user_id != +req.params.user_id){
		// not an admin and not your own user
		return res.status(403).json({ message: "change not allowed" });
	}

	if (req.claims.role === "UPLOAD_ONLY"){
		// static account, no changes happen here
		return res.status(403).json({ message: "change not allowed (static account)" });
	}

	if (!req.body.name){
		return res.status(400).json({ message: "no new name specified" });
	}

	if (req.body.name !== req.claims.username){
		// check if new name is available
		const existing_user = await user_dao.get_by_name(req.body.name);
		if (existing_user){
			return res.status(400).json({ message: "user with that name already exists" });
		}
	}

	const user = await user_dao.get_by_id(+req.params.user_id);
	if (!user){
		return res.status(400).json({ message: "user does not exist" });
	}

	logger.info(`user '${req.claims.username}' with id ${req.claims.user_id} changes name to '${req.body.name}' and mail to '${req.body.mail}'`);
	user.name = req.body.name;
	user.mail = req.body.mail;
	await user_dao.update_user(user);

	res.json({
		success: true
	});
});
