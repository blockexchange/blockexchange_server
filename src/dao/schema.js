const pool = require("../pool");

module.exports.create = function(data) {
  const query = `
    insert into
    schema(
      complete, user_id, uid, description,
      size_x, size_y, size_z, part_length,
      total_size, total_parts, created
    )
    values(
      $1, $2, $3, $4,
      $5, $6, $7, $8,
      $9, $10, $11
    )
    returning *
  `;

  const values = [
    false, data.user_id, data.uid, data.description || "",
    data.size_x, data.size_y, data.size_z, data.part_length,
    0, 0, Date.now()
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

module.exports.finalize = function(schema_id) {
  const query = `
    update schema s
    set complete = true,
		total_size = (select sum(length(data)) + sum(length(metadata)) from schemapart sp where sp.schema_id = s.id),
		total_parts = (select count(*) from schemapart sp where sp.schema_id = s.id)
    where id = $1
  `;

  return new Promise(function(resolve, reject){
    pool.connect()
    .then(client => {
      client.query(query, [schema_id])
      .then(() => {
  			client.release();
  			resolve();
  		})
      .catch(e => {
  			client.release();
        console.error(e.stack);
        reject();
      });
    });
  });
};

module.exports.get_by_uid = function(uid) {
  return new Promise(function(resolve, reject) {
    pool.connect()
    .then(client => {
      client.query("select * from schema where uid = $1", [uid])
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
