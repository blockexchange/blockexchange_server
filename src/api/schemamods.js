const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schemamod_dao = require("../dao/schemamod");
const schema_dao = require("../dao/schema");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");

const { verifytoken } = require("../middleware/token");

const schemagroup_verify = function(req, res, next){
  const user_id = req.claims.user_id;
  const schema_id = req.params.id;

  schema_dao.get_by_id(schema_id)
  .then(schema => user_schemagroup_permission_dao.get_by_user_and_schemagroup_id(user_id, schema.schemagroup_id))
  .then(perm => {
    if (!perm || !perm.delete){
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

app.get('/api/schema/:id/mods', function(req, res){
  console.log("GET /api/schema/:id/mods", req.params.id);

  schema_dao.get_by_id(req.params.id)
  .then(schema => {
    return schemamod_dao.find_all(schema.id)
    .then(schemamods => {
      const result = {};
      schemamods.forEach(sm => {
        result[sm.mod_name] = sm.node_count;
      });
      res.json(result);
    });
  })
  .catch(() => res.status(500).end());

});


app.post('/api/schema/:id/mods', verifytoken, schemagroup_verify, jsonParser, function(req, res){
  console.log("POST /api/schema/:id/mods", req.params.id, req.body);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

    Object.keys(req.body).forEach(mod_name => {
      const node_count = req.body[mod_name];
      schemamod_dao.create(schema.id, mod_name, node_count);
    });
    res.end();
  });
});
