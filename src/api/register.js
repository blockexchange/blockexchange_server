const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const bcrypt = require('bcryptjs');

const app = require("../app");
const user_dao = require("../dao/user");

// data='{"name": "xyz", "password": "abc", "mail": null}'
// curl -X POST 127.0.0.1:8080/api/register --data "${data}" -H "Content-Type: application/json"
app.post('/api/register', jsonParser, function(req, res){
  console.log("POST /api/register", req.body.name, req.body.mail);

  if (!req.body.name || req.body.name.length == 0) {
    res.json({
      success: false,
      message: "Invalid username"
    });
    return;
  }

  if (!req.body.password || req.body.password.length == 0) {
    res.json({
      success: false,
      message: "Invalid password"
    });
    return;
  }

  user_dao.get_by_name(req.body.name)
  .then(user => {
    if (user) {
      res.json({
        success: false,
        message: "Username already exists"
      });
    } else {
      var salt = bcrypt.genSaltSync(10);
      var hash = bcrypt.hashSync(req.body.password, salt);
      return user_dao.create({
          role: "MEMBER",
          name: req.body.name,
          hash: hash,
          mail: req.body.mail
      })
      .then(() => res.json({ success: true }));
    }
  })
  .catch((e) => {
    console.error(e);
    res.status(500).end();
  });
});
