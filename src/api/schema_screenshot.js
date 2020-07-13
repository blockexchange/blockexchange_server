const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");
const schema_screenshot_dao = require("../dao/schema_screenshot");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");

const { verifytoken, rolecheck } = require("../middleware/token");

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


app.get('/api/schema/:id/screenshot', function(req, res){
  console.log("GET /api/schema/:id/screenshot", req.params.id);

  schema_screenshot_dao.find_all(req.params.id)
  .then(screenshots => screenshots || [])
  .then(screenshots => res.json(screenshots))
  .catch(() => res.status(500).end());
});

app.get('/api/schema/:id/screenshot/:screenshot_id', function(req, res){
  console.log("GET /api/schema/:id/screenshot/:screenshot_id", req.params.id, req.params.screenshot_id);

  schema_screenshot_dao.get_by_id(req.params.screenshot_id)
  .then(screenshot => {
    res.header("")
    .send(screenshot.data);
  })
  .catch(() => res.status(500).end());
});

app.post('/api/schema/:id/screenshot', verifytoken, rolecheck("MEMBER"), schemagroup_verify, jsonParser, function(req, res){
  console.log("POST /api/schema/:id/screenshot", req.params.id);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

    schema_screenshot_dao.create(schema.id, req.body.title, req.body.type, req.body.data)
    .then(() => res.end());
  });
});


app.get('/api/schema/:id/screenshot/:screenshot_id', verifytoken, rolecheck("MEMBER"), schemagroup_verify, function(req, res){
  console.log("DELETE /api/schema/:id/screenshot/:screenshot_id", req.params.id, req.params.screenshot_id);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    // check user id in claims
    if (schema.user_id != +req.claims.user_id){
      res.status(401).end();
      return;
    }

    schema_screenshot_dao.remove(req.params.screenshot_id)
    .then(() => res.end());
  })
  .catch(() => res.status(500).end());
});
