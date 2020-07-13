const bodyParser = require('body-parser');
const jsonParser = bodyParser.json({limit: '20mb'});

const app = require("../app");
const serializer = require("../util/serializer");

const schemapart_dao = require("../dao/schemapart");
const schema_dao = require("../dao/schema");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");

const { verifytoken } = require("../middleware/token");

const schemagroup_verify = function(req, res, next){
  const user_id = req.claims.user_id;
  const schema_id = req.params.schema_id;

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

app.post('/api/schemapart', verifytoken, schemagroup_verify, jsonParser, function(req, res){
  console.log("POST /api/schemapart", req.body.schema_id, req.body.offset_x, req.body.offset_y, req.body.offset_z);

  return schema_dao.get_by_id(req.body.schema_id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

    if (schema.completed) {
      res.status(500).end();
      return;
    }

    const serialized_data = serializer.serialize(req.body.data);

    schemapart_dao.create({
      schema_id: schema.id,
      offset_x: req.body.offset_x,
      offset_y: req.body.offset_y,
      offset_z: req.body.offset_z,
      data: serialized_data.data,
      metadata: serialized_data.metadata
    })
    .then(id_obj => res.json(id_obj))
    .catch(() => res.status(500).end());
  });
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
    if (schemapart) {
      const data = serializer.deserialize(schemapart);

      res.json({
        schema_id: req.params.schema_id,
        offset_x: req.params.offset_x,
        offset_y: req.params.offset_y,
        offset_z: req.params.offset_z,
        data: {
          node_ids: data.node_ids,
          param1: data.param1,
          param2: data.param2,
          metadata: data.metadata,
          size: data.size,
					node_mapping: data.node_mapping
        }
      });
    } else
      res.status(404).end();
  })
  .catch(e => {
		console.error(e);
		res.status(500).end();
	});
});
