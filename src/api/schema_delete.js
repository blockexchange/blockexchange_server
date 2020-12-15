const app = require("../app");
const schema_dao = require("../dao/schema");
const logger = require("../logger");

const tokenmiddleware = require("../middleware/token");
const rolecheck = require("../util/rolecheck");
const tokencheck = tokenmiddleware(claims => rolecheck.can_delete(claims.role));


app.delete("/api/schema/:id", tokencheck, function(req, res){
  logger.debug("DELETE /api/schema/:id", req.params.id, req.body);

  return schema_dao.get_by_id(req.params.id)
  .then(schema => {
    if (schema.user_id != req.claims.user_id){
      res.status(403).end();
      return;
    }

    // delete schema
    return schema_dao.delete_by_id(req.params.id)
    .then(() => res.end());
  });
});
