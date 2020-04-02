const pool = require("../pool");
const zlib = require('zlib');

module.exports.create = function(schemapart) {
  const query = `
    insert into
    schemapart(schema_id, offset_x, offset_y, offset_z, data)
    values($1, $2, $3, $4, $5)
    returning id
  `;

  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {
  		//TODO: check if schema is not completed
  	  const compressed = zlib.gzipSync(schemapart.data);

  	  const values = [
  	    schemapart.schema_id,
  	    schemapart.offset_x,
  	    schemapart.offset_y,
  	    schemapart.offset_z,
  	    compressed
  	  ];

      client.query(query, values)
      .then(sql_res => {
  			resolve(sql_res.rows[0]);
  			client.release();
  		})
      .catch(e => {
  			client.release();
        console.error(e.stack);
        reject();
      });
    });
  });
};

module.exports.get_by_id_and_offset = function(schema_id, x, y, z) {
  const query = `
    select *
    from schemapart
    where schema_id = $1
    and offset_x = $2
    and offset_y = $3
    and offset_z = $4
  `;

  return new Promise(function(resolve, reject){
    const values = [ schema_id, x, y, z ];

    pool.connect()
    .then(client => {
      client.query(query, values)
      .then(sql_res => {
        if (sql_res.rowCount == 0){
          resolve();
        } else {
          const row = sql_res.rows[0];
          const part = Object.assign({}, row);
          part.data = zlib.gunzipSync(row.data).toString("utf-8");
          resolve(part);
        }
  			client.release();
      })
      .catch(e => {
  			client.release();
        console.error(e.stack);
        reject();
      });
    });
  });
};
