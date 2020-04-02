const bodyParser = require('body-parser');
const jsonParser = bodyParser.json({limit: '20mb'});

const app = require("../app");
const schemapart_dao = require("../dao/schemapart");

// data='{"schema_id": 1, "offset_x": 0, "offset_y": 0, "offset_z": 0, "data": "return {}"}'
// curl -X POST 127.0.0.1:8080/api/schemapart --data "${data}" -H "Content-Type: application/json"
app.post('/api/schemapart', jsonParser, function(req, res){
  console.log("POST /api/schemapart", req.body.schema_id, req.body.offset_x, req.body.offset_y, req.body.offset_z);
	//console.log("Data: ", req.body);//DEBUG

  schemapart_dao.create({
    schema_id: req.body.schema_id,
    offset_x: req.body.offset_x,
    offset_y: req.body.offset_y,
    offset_z: req.body.offset_z,
    data: req.body.data
  })
  .then(id_obj => res.json(id_obj))
  .catch(() => res.status(500).end());
});


// curl 127.0.0.1:8080/api/schemapart/1/0/0/0
app.get('/api/schemapart/:schema_id/:offset_x/:offset_y/:offset_z', function(req, res){
  console.log("GET /api/schemapart", req.params);

  schemapart_dao.get_by_id_and_offset(
    req.params.schema_id,
    req.params.offset_x,
    req.params.offset_y,
    req.params.offset_z
  )
  .then(schemapart => {
    if (schemapart)
      res.json(schemapart);
    else
      res.status(404).end();
  })
  .catch(() => res.status(500).end());
});
