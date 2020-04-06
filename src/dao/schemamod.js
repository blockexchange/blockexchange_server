const pool = require("../pool");

module.exports.create = function(schema_id, mod_name, node_count) {
  const query = `
    insert into
    schemamod(schema_id, mod_name, node_count)
    values($1, $2, $3)
    returning *
  `;

  const values = [
    schema_id, mod_name, node_count
  ];

  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {
      return client.query(query, values)
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

module.exports.find_all = function(schema_id) {
  const query = `
    select * from schemamod
    where schema_id = $1
  `;

  const values = [
    schema_id
  ];

  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {
      return client.query(query, values)
      .then(sql_res => {
        resolve(sql_res.rows);
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
