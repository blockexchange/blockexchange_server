const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const schema_dao = require("../dao/schema");

// data='{"keywords": "description"}'
// curl -X POST 127.0.0.1:8080/api/searchschema --data "${data}" -H "Content-Type: application/json"
app.post('/api/searchschema', jsonParser, function(req, res){
  console.log("POST /api/searchschema", req.body);

	var q;
	if (req.body.user_id) {
		q = schema_dao.find_by_user_id(req.body.user_id);
	} else if (req.body.keywords) {
		q = schema_dao.find_by_description(req.body.keywords);
	}

  q.then(rows => res.json(rows))
  .catch(e => {
    console.error(e);
    res.status(500).end();
  });

});

app.get("/api/searchrecent/:count", function(req, res){
	schema_dao.find_recent(req.params.count)
	.then(rows => res.json(rows))
	.catch(e => {
		console.error(e);
		res.status(500).end();
	});

});
