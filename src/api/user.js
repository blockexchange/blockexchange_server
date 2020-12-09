const app = require("../app");
const user_dao = require("../dao/user");
const logger = require("../logger");


app.get('/api/user', function(req, res){
  logger.debug("GET /api/user");

  user_dao.get_all()
  .then(users => {
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
  })
  .catch(() => res.status(500).end());

});
