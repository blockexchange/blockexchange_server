const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");
const { verifytoken, rolecheck } = require("../middleware/token");

const schemagroup_verify = function(req, res, next){
  const schemagroup_id = req.body.schemagroup_id;
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

app.post('/api/schema', verifytoken, rolecheck("MEMBER"), jsonParser, schemagroup_verify, function(req, res){
  console.log("POST /api/schema", req.body);

  schema_dao.create({
    schemagroup_id: req.body.schemagroup_id,
    name: req.body.name,
    description: req.body.description,
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
