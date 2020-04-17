
const app = require("../app");
const schema_dao = require("../dao/schema");

// curl 127.0.0.1:8080/api/schema/1
app.get('/api/schema/:id', function(req, res){
  console.log("GET /api/schema/:id", req.params.id);

  schema_dao.get_by_id(req.params.id)
  .then(schema => res.json(schema))
  .catch(() => res.status(500).end());
});
