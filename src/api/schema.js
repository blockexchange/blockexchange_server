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
    if (!claims.permissions.schema.create){
      res.status(401).end();
      return;
    }

    schema_dao.create({
      user_id: +claims.user_id,
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
  })
  .catch(e => {
    console.error(e);
    res.status(401).end();
  });

});

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


// curl -X POST 127.0.0.1:8080/api/schema/1/complete
app.post('/api/schema/:id/complete', jsonParser, function(req, res){
  console.log("POST /api/schema/id/complete", req.params.id, req.body);

  tokencheck(req, res)
  .then(claims => {
    // check if the user can create schemas
    if (!claims.permissions.schema.create){
      res.status(401).end();
      return;
    }

    return schema_dao.get_by_id(req.params.id)
    .then(schema => {
      // check user id in claims
      if (schema.user_id != +claims.user_id){
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
  })
  .catch(() => res.status(401).end());
});


// curl 127.0.0.1:8080/api/schema/1
app.get('/api/schema/:id', function(req, res){
  console.log("GET /api/schema/:id", req.params.id);

  schema_dao.get_by_id(req.params.id)
  .then(schema => res.json(schema))
  .catch(() => res.status(500).end());
});
