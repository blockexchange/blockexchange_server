const executor = require("./executor");

module.exports.create = function(user_id, schema_id) {
  const query = `
    insert into
    user_schema_star(user_id, schema_id)
    values($1, $2)
    returning *
  `;

  const values = [
    user_id, schema_id
  ];

  return executor(query, values, { single_row: true });
};

module.exports.remove = function(user_id, schema_id) {
  const query = `
    delete from
    user_schema_star
		where user_id = $1 and
		schema_id = $2
  `;

  const values = [
    user_id, schema_id
  ];

  return executor(query, values);
};

module.exports.find_by_schema_id = function(schema_id) {
  const query = `
    select * from user_schema_star
    where schema_id = $1
  `;

  const values = [
    schema_id
  ];

  return executor(query, values, { single_row: true });
};


module.exports.count_by_schema_id = function(schema_id) {
  const query = `
    select count(*) as stars from user_schema_star
    where schema_id = $1
  `;

  const values = [
    schema_id
  ];

  return executor(query, values, { single_row: true });
};
