const app = require("../app");
const schema_dao = require("../dao/schema");
const tokencheck = require("../util/tokencheck");


app.delete("/api/schema/:id", function(req, res){
  console.log("DELETE /api/schema/:id", req.params.id, req.body);

  tokencheck(req, res)
  .then(claims => {
    if (!claims.permissions.schema.delete){
      res.status(401).end();
      return;
    }

    return schema_dao.get_by_id(req.params.id)
    .then(schema => {
      if (schema.user_id != claims.user_id){
        res.status(401).end();
        return;
      }

      // delete schema
      return schema_dao.delete_by_id(req.params.id)
      .then(() => res.end());
    });
  })
  .catch(e => {
    console.error(e);
    res.status(500).end();
  });
});
