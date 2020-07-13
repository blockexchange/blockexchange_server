const app = require("../app");

const schema_dao = require("../dao/schema");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");

const { verifytoken, rolecheck } = require("../middleware/token");

const schemagroup_verify = function(req, res, next){
  const user_id = req.claims.user_id;
  const schema_id = req.params.id;

  schema_dao.get_by_id(schema_id)
  .then(schema => user_schemagroup_permission_dao.get_by_user_and_schemagroup_id(user_id, schema.schemagroup_id))
  .then(perm => {
    if (!perm || !perm.update){
      res.status(403).end();
      return;
    }
    next();
  })
  .catch(e => {
    console.error(e);
    res.status(500).end();
  });
};


app.delete("/api/schema/:id", verifytoken, rolecheck("MEMBER"), schemagroup_verify, function(req, res){
  console.log("DELETE /api/schema/:id", req.params.id, req.body);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    if (schema.user_id != req.claims.user_id){
      res.status(403).end();
      return;
    }

    // delete schema
    return schema_dao.delete_by_id(req.params.id)
    .then(() => res.end());
  });
});
