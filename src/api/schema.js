const bodyParser = require('body-parser');
const jsonParser = bodyParser.json();

const app = require("../app");
const pool = require("../pool");

// data='{"size_x": 10, "size_y": 10, "size_z": 10, "part_length": 10}'
// curl -X POST 127.0.0.1:8080/api/schema --data "${data}" -H "Content-Type: application/json"
app.post('/api/schema', jsonParser, function(req, res){
  console.log("POST /api/schema", req.body);

  const query = `
    insert into
    schema(complete, size_x, size_y, size_z, part_length)
    values($1, $2, $3, $4, $5)
    returning *
  `;

  const values = [
    false,
    req.body.size_x,
    req.body.size_y,
    req.body.size_z,
    req.body.part_length
  ];

  pool.connect()
  .then(client => {
    client.query(query, values)
    .then(sql_res => {
			res.json(sql_res.rows[0]);
			client.release();
		})
    .catch(e => {
			client.release();
      console.error(e.stack);
      res.status(500).end();
    });
  });

});

// curl -X POST 127.0.0.1:8080/api/schema/1/complete
app.post('/api/schema/:id/complete', function(req, res){
  console.log("POST /api/schema/id/complete", req.params.id);

  const query = `
    update schema
    set complete = true
    where id = $1
  `;

  pool.connect()
  .then(client => {
    client.query(query, [req.params.id])
    .then(() => {
			client.release();
			res.end();
		})
    .catch(e => {
			client.release();
      console.error(e.stack);
      res.status(500).end();
    });
  });

});


// curl 127.0.0.1:8080/api/schema/1
app.get('/api/schema/:id', function(req, res){
  console.log("GET /api/schema", req.params.id);

  pool.connect()
  .then(client => {
    client.query("select * from schema where id = $1", [req.params.id])
    .then(sql_res => {
			res.json(sql_res.rows[0]);
			client.release();
		})
    .catch(e => {
			client.release();
      console.error(e.stack);
      res.status(500).end();
    });
  });
});
