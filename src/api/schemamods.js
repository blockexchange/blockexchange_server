const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schemamod_dao = require("../dao/schemamod");
const schema_dao = require("../dao/schema");
const tokencheck = require("../util/tokencheck");


app.get('/api/schema/:uid/mods', function(req, res){
  console.log("GET /api/schema/:uid/mods", req.params.uid);

  schema_dao.get_by_uid(req.params.uid)
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


app.post('/api/schema/:uid/mods', jsonParser, function(req, res){
  console.log("POST /api/schema/:id/mods", req.params.id, req.body);

  tokencheck(req, res)
  .then(claims => {
    return schema_dao.get_by_uid(req.params.uid)
    .then(schema => {
      // check user id in claims
      if (schema.user_id != +claims.user_id){
        res.status(401).end();
        return;
      }

      Object.keys(req.body).forEach(mod_name => {
        const node_count = req.body[mod_name];
        schemamod_dao.create(schema.id, mod_name, node_count);
      });
      res.end();
    });
  })
  .catch(() => res.status(401).end());
});
