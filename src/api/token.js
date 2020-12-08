const bcrypt = require('bcryptjs');
const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const createjwt = require("../util/createjwt");
const logger = require("../logger");
const app = require("../app");
const user_dao = require("../dao/user");

// data='{"name": "xyz", "password": "abc"}'
// curl -X POST 127.0.0.1:8080/api/token --data "${data}" -H "Content-Type: application/json"
app.post('/api/token', jsonParser, function(req, res){
  logger.debug("POST /api/token", req.body.name);

  user_dao.get_by_name(req.body.name)
  .then(user => {
    if (!user) {
      res.status(404).json({
				message: "User not found"
			});
      return;
    }

		if (user.type != "LOCAL"){
			res.status(403).json({
				message: "Direct login with an external user not allowed"
			});
      return;
		}

    const success = bcrypt.compareSync(req.body.password, user.hash);
    if (!success) {
      res.status(401).json({
				message: "Invalid password"
			});
      return;
    }

		const token = createjwt(user);
    res.send(token);
  })
  .catch((e) => {
    console.error(e);
    res.status(500).end({
			message: "Internal server error"
		});
  });
});
