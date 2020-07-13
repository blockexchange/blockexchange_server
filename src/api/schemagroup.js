
const app = require("../app");
const schemagroup_dao = require("../dao/schemagroup");
const user_schemagroup_permission_dao = require("../dao/user_schemagroup_permission");

app.get('/api/schemagroup/:id', function(req, res){
  console.log("GET /api/schemagroup/:id", req.params.id);

  Promise.all([
    user_schemagroup_permission_dao.get_by_schemagroup_id(req.params.id),
    schemagroup_dao.get_by_id(req.params.id)
  ])
  .then(([permissions, schemagroup]) => {
    const result = Object.assign({}, schemagroup, {
      permissions: permissions
    });
    res.json(result);
  })
  .catch(() => res.status(500).end());
});
