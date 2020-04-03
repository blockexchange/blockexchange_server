const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const bcrypt = require('bcrypt');
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

    bcrypt.compare(req.body.password, user.hash)
    .then(success => {
      if (!success) {
        res.status(401).end();
      }

      const payload = {
        username: user.name,
        userid: user.id
      };
      const token = jwt.sign(payload, process.env.BLOCKEXCHANGE_KEY);
      res.send(token);
    });
  })
  .catch((e) => {
    console.error(e);
    res.status(500).end();
  });
});
