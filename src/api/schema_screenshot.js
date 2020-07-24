const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_screenshot_dao = require("../dao/schema_screenshot");
const schema_dao = require("../dao/schema");

const tokenmiddleware = require("../middleware/token");
const permission_create = tokenmiddleware(claims => claims.permissions.screenshot.create);
const permission_delete = tokenmiddleware(claims => claims.permissions.screenshot.delete);

app.get('/api/schema/:id/screenshot', function(req, res){
  console.log("GET /api/schema/:id/screenshot", req.params.id);

  schema_screenshot_dao.find_all(req.params.id)
  .then(screenshots => screenshots || [])
  .then(screenshots => screenshots.map(s => ({ id: s.id, type: s.type, title: s.title })))
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

app.post('/api/schema/:id/screenshot', permission_create, jsonParser, function(req, res){
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


app.get('/api/schema/:id/screenshot/:screenshot_id', permission_delete, function(req, res){
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
