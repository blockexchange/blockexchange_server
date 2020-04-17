const executor = require("./executor");

module.exports.create = function(schemapart) {
  const query = `
    insert into
    schemapart(schema_id, offset_x, offset_y, offset_z, data, metadata)
    values($1, $2, $3, $4, $5, $6)
    returning id
  `;

  const values = [
    schemapart.schema_id,
    schemapart.offset_x,
    schemapart.offset_y,
    schemapart.offset_z,
    schemapart.data,
    schemapart.metadata
  ];

  return executor(query, values, { single_row: true });
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

  const values = [ schema_id, x, y, z ];

  return executor(query, values, { single_row: true });
};
