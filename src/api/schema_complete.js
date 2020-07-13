const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");
const { verifytoken } = require("../middleware/token");

const schemagroup_verify = function(req, res, next){
  const schemagroup_id = req.params.id;
  const user_id = req.claims.user_id;

  user_schemagroup_permission_dao.get_by_user_and_schemagroup_id(user_id, schemagroup_id)
  .then(perm => {
    if (!perm || !perm.create){
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


app.post('/api/schema/:id/complete', verifytoken, jsonParser, schemagroup_verify, function(req, res){
  console.log("POST /api/schema/id/complete", req.params.id, req.body);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check if already completed
    if (schema.complete){
      res.status(500).end();
      return;
    }

    return schema_dao.finalize(schema.id)
    .then(() => res.end())
    .catch(() => res.status(500).end());
    });
});
