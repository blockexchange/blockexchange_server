const pool = require("../pool");

module.exports.create = function(data) {
  const query = `
    insert into
    schema(
      complete, user_id, name, description,
      size_x, size_y, size_z, part_length,
      total_size, total_parts, created,
      search_tokens
    )
    values(
      $1, $2, $3, $4,
      $5, $6, $7, $8,
      $9, $10, $11,
      to_tsvector($12)
    )
    returning *
  `;

  const values = [
    false, data.user_id, data.name, data.description || "",
    data.size_x, data.size_y, data.size_z, data.part_length,
    0, 0, Date.now(),
    data.name + " " + data.description
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

module.exports.get_by_name = function(name) {
  return new Promise(function(resolve, reject) {
    pool.connect()
    .then(client => {
      client.query("select * from schema where name = $1", [name])
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

module.exports.get_by_id = function(id) {
  return new Promise(function(resolve, reject) {
    pool.connect()
    .then(client => {
      client.query("select * from schema where id = $1", [id])
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

module.exports.find_by_user_id = function(user_id) {
  return new Promise(function(resolve, reject) {
    pool.connect()
    .then(client => {
      client.query("select * from schema where user_id = $1", [user_id])
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

module.exports.find_by_keywords = function(keywords) {
  const query = `
    select *
    from schema
    where search_tokens @@ to_tsquery($1)
    limit 1000
  `;
  return new Promise(function(resolve, reject) {
    pool.connect()
    .then(client => {
      client.query(query, [keywords])
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


module.exports.find_recent = function(count) {
  const query = `
    select *
    from schema
		order by created desc
    limit $1
  `;
  return new Promise(function(resolve, reject) {
    pool.connect()
    .then(client => {
      client.query(query, [Math.min(count, 250)])
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
