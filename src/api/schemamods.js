const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schemamod_dao = require("../dao/schemamod");
const schema_dao = require("../dao/schema");

const tokenmiddleware = require("../middleware/token");
const tokencheck = tokenmiddleware(claims => {
  return claims.permissions.schema.create;
});

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


app.post('/api/schema/:id/mods', tokencheck, jsonParser, function(req, res){
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
