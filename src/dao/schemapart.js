const pool = require("../pool");

module.exports.create = function(schemapart) {
  const query = `
    insert into
    schemapart(schema_id, offset_x, offset_y, offset_z, data, metadata)
    values($1, $2, $3, $4, $5, $6)
    returning id
  `;

  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {

  	  const values = [
  	    schemapart.schema_id,
  	    schemapart.offset_x,
  	    schemapart.offset_y,
  	    schemapart.offset_z,
  	    schemapart.data,
        schemapart.metadata
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
          resolve(sql_res.rows[0]);
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
