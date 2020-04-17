const executor = require("./executor");

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

  return executor(query, values, { single_row: true });
};

module.exports.find_all = function(schema_id) {
  const query = `
    select * from schemamod
    where schema_id = $1
  `;

  const values = [
    schema_id
  ];

  return executor(query, values);
};
