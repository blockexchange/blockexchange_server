const executor = require("./executor");

module.exports.create = function(schema_id, tag_id) {
  const query = `
    insert into
    schematag(schema_id, tag_id)
    values($1, $2)
    returning *
  `;

  const values = [
    schema_id, tag_id
  ];

  return executor(query, values, { single_row: true });
};

module.exports.find = function(schema_id) {
  const query = `
    select * from schematag
    where schema_id = $1
  `;

  const values = [
    schema_id
  ];

  return executor(query, values);
};

module.exports.remove = function(schema_id, tag_id) {
  const query = `
    delete from schematag
    where schema_id = $1
		and tag_id = $2
  `;

  const values = [
    schema_id, tag_id
  ];

  return executor(query, values);
};
