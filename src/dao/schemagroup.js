const executor = require("./executor");

module.exports.get_by_name = function(name) {
  return executor("select * from schemagroup where name = $1", [name], { single_row: true });
};

module.exports.get_by_id = function(id) {
  return executor("select * from schemagroup where id = $1", [id], { single_row: true });
};

module.exports.get_all = function() {
  return executor("select * from schemagroup");
};

module.exports.create = function(data) {
  const query = `insert into schemagroup(name, created) values($1, $2) returning *`;
  const values = [data.name, Date.now()];

  return executor(query, values, { single_row: true });
};

module.exports.remove_by_name = function(name) {
  const query = `delete from schemagroup where name = $1`;
  const values = [name];

  return executor(query, values);
};
