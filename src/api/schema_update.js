const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const tokencheck = require("../util/tokencheck");

app.put("/api/schema/:id", jsonParser, function(req, res){
  console.log("PUT /api/schema/:id", req.params.id, req.body);
  const new_schema = req.body;

  tokencheck(req, res)
  .then(claims => {
    if (!claims.permissions.schema.update){
      res.status(401).end();
      return;
    }

    return schema_dao.get_by_id(new_schema.id)
    .then(old_schema => {
      if (old_schema.user_id != claims.user_id){
        res.status(401).end();
        return;
      }

      // update schema
      return schema_dao.update(new_schema)
      .then(() => res.json(new_schema));
    });
  })
  .catch(e => {
    console.error(e);
    res.status(500).end();
  });
});
