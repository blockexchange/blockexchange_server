const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const userschemastar_dao = require("../dao/userschemastar");
const tokenmiddleware = require("../middleware/token");
const tokencheck = tokenmiddleware();

app.post('/api/schema/:id/star', tokencheck, jsonParser, function(req, res){
  console.log("POST /api/schema/:id/star", req.params.id, req.claims.user_id);

	userschemastar_dao.create(req.claims.user_id, req.params.id)
	.then(() => res.end())
	.catch(() => res.status(500).end());

});

app.delete('/api/schema/:id/star', tokencheck, jsonParser, function(req, res){
  console.log("DELETE /api/schema/:id/star", req.params.id, req.claims.user_id);

	userschemastar_dao.remove(req.claims.user_id, req.params.id)
	.then(() => res.end())
	.catch(() => res.status(500).end());

});
