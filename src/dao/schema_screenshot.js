const executor = require("./executor");

module.exports.create = function(schema_id, title, type, data) {
  const query = `
    insert into
    schema_screenshot(schema_id, title, type, data)
    values($1, $2, $3, $4)
    returning *
  `;

  const values = [schema_id, title, type, data];

  return executor(query, values, { single_row: true });
};

module.exports.find_all = function(schema_id) {
  const query = `select id, type, title from schema_screenshot where schema_id = $1`;
  const values = [schema_id];

  return executor(query, values);
};

module.exports.get_by_id = function(schema_id, id) {
  const query = `select * from schema_screenshot where schema_id = $1 and id = $2`;
  const values = [schema_id, id];

  return executor(query, values, { single_row: true });
};


module.exports.remove = function(id) {
  const query = `delete from schema_screenshot where id = $1`;
  const values = [id];

  return executor(query, values);
};
