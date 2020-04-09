const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const tokencheck = require("../util/tokencheck");

// data='{"size_x": 10, "size_y": 10, "size_z": 10, "part_length": 10, "name": "xyz"}'
// curl -X POST 127.0.0.1:8080/api/schema --data "${data}" -H "Content-Type: application/json"
app.post('/api/schema', jsonParser, function(req, res){
  console.log("POST /api/schema", req.body);

  tokencheck(req, res)
  .then(claims => {
    schema_dao.create({
      user_id: +claims.user_id,
      name: req.body.name,
      description: req.body.description,
      size_x: req.body.size_x,
      size_y: req.body.size_y,
      size_z: req.body.size_z,
      part_length: req.body.part_length
    })
    .then(inserted_data => res.json(inserted_data))
    .catch(e => {
      console.error(e);
      res.status(500).end();
    });
  })
  .catch(e => {
    console.error(e);
    res.status(401).end();
  });

});

// curl -X POST 127.0.0.1:8080/api/schema/1/complete
app.post('/api/schema/:id/complete', jsonParser, function(req, res){
  console.log("POST /api/schema/id/complete", req.params.id, req.body);

  tokencheck(req, res)
  .then(claims => {
    return schema_dao.get_by_id(req.params.id)
    .then(schema => {
      console.log(schema, claims);
      // check user id in claims
      if (schema.user_id != +claims.user_id){
        res.status(401).end();
        return;
      }

      return schema_dao.finalize(schema.id)
      .then(() => res.end())
      .catch(() => res.status(500).end());
    });
  })
  .catch(() => res.status(401).end());
});


// curl 127.0.0.1:8080/api/schema/1
app.get('/api/schema/:id', function(req, res){
  console.log("GET /api/schema/:id", req.params.id);

  schema_dao.get_by_id(req.params.id)
  .then(schema => {
    res.json(schema);
    schema.downloads++;
    schema_dao.update(schema);
  })
  .catch(() => res.status(500).end());
});
