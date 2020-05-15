const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const bcrypt = require('bcryptjs');
var jwt = require('jsonwebtoken');

const app = require("../app");
const user_dao = require("../dao/user");

// data='{"name": "xyz", "password": "abc"}'
// curl -X POST 127.0.0.1:8080/api/token --data "${data}" -H "Content-Type: application/json"
app.post('/api/token', jsonParser, function(req, res){
  console.log("POST /api/token", req.body.name);

  user_dao.get_by_name(req.body.name)
  .then(user => {
    if (!user) {
      res.status(404).end();
      return;
    }

    const success = bcrypt.compareSync(req.body.password, user.hash);
    if (!success) {
      res.status(401).end();
      return;
    }

    const payload = {
      username: user.name,
      user_id: user.id
    };

    // check for temporary role, only allow creation of content with it
    if (user.role == "TEMP"){
      // temporary/default user
      payload.permissions = {
        schema: {
          create: true
        }
      };
    } else {
      // normal user
      payload.permissions = {
        user: {
          update: true
        },
        schema: {
          create: true,
          update: true,
          delete: true
        },
        screenshot: {
          create: true,
          delete: true
        }
      };
    }
    // TODO: custom permissions on token for normal users

    const token = jwt.sign(payload, process.env.BLOCKEXCHANGE_KEY);
    res.send(token);
  })
  .catch((e) => {
    console.error(e);
    res.status(500).end();
  });
});
