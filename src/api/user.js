const app = require("../app");
const user_dao = require("../dao/user");


app.get('/api/user', function(req, res){
  console.log("GET /api/user");

  user_dao.get_all()
  .then(users => {
    const list = users.map(user => {
      return {
        name: user.name,
        id: user.id,
        created: user.created
      };
    });

    res.json(list);
  })
  .catch(() => res.status(500).end());

});
