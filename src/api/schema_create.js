const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const tokenmiddleware = require("../middleware/token");
const tokencheck = tokenmiddleware(claims => {
  return claims.permissions.schema.create;
});

app.post('/api/schema', tokencheck, jsonParser, function(req, res){
  console.log("POST /api/schema", req.body);

  schema_dao.create({
    user_id: +req.claims.user_id,
    name: req.body.name,
    description: req.body.description,
    long_description: req.body.long_description,
    size_x: req.body.size_x,
    size_y: req.body.size_y,
    size_z: req.body.size_z,
    part_length: req.body.part_length,
		license: req.body.license
  })
  .then(inserted_data => res.json(inserted_data))
  .catch(e => {
    console.error(e);
    res.status(500).end();
  });

});



app.post('/api/schema/:id/complete', tokencheck, jsonParser, function(req, res){
  console.log("POST /api/schema/id/complete", req.params.id, req.body);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

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
