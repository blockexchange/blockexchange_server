const bodyParser = require('body-parser');
const zlib = require('zlib');
const jsonParser = bodyParser.json();

const app = require("../app");
const pool = require("../pool");

// data='{"schema_id": 1, "offset_x": 0, "offset_y": 0, "offset_z": 0, "data": "return {}"}'
// curl -X POST 127.0.0.1:8080/api/schemapart --data "${data}" -H "Content-Type: application/json"
app.post('/api/schemapart', jsonParser, function(req, res){
  console.log("POST /api/schemapart", req.body.schema_id, req.body.offset_x, req.body.offset_y, req.body.offset_z);

  const query = `
    insert into
    schemapart(schema_id, offset_x, offset_y, offset_z, data)
    values($1, $2, $3, $4, $5)
    returning id
  `;

  const compressed = zlib.gzipSync(req.body.data);

  const values = [
    req.body.schema_id,
    req.body.offset_x,
    req.body.offset_y,
    req.body.offset_z,
    compressed
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


// curl 127.0.0.1:8080/api/schemapart/1/0/0/0
app.get('/api/schemapart/:schema_id/:offset_x/:offset_y/:offset_z', function(req, res){
  console.log("GET /api/schemapart", req.params);

  const query = `
    select *
    from schemapart
    where schema_id = $1
    and offset_x = $2
    and offset_y = $3
    and offset_z = $4
  `;

  const values = [
    req.params.schema_id,
    req.params.offset_x,
    req.params.offset_y,
    req.params.offset_z
  ];

  pool.connect()
  .then(client => {
    client.query(query, values)
    .then(sql_res => {
      if (sql_res.rowCount == 0){
        res.status(404).end();
      } else {
        const row = sql_res.rows[0];
        const part = Object.assign({}, row);
        part.data = zlib.gunzipSync(row.data).toString("utf-8");
        res.json(part);
      }
			client.release();
    })
    .catch(e => {
			client.release();
      console.error(e.stack);
      res.status(500).end();
    });
  });
});
