const executor = require("./executor");

module.exports.create = function(schema_id, title, data) {
  const query = `
    insert into
    schema_screenshot(schema_id, title, data)
    values($1, $2, $3)
    returning *
  `;

  const values = [
    schema_id, title, data
  ];

  return executor(query, values, { single_row: true });
};

module.exports.find_all = function(schema_id) {
  const query = `
    select * from schema_screenshot
    where schema_id = $1
  `;

  const values = [
    schema_id
  ];

  return executor(query, values);
};


module.exports.remove = function(id) {
  const query = `delete from schema_screenshot where id = $1`;
  const values = [id];

  return executor(query, values);
};
