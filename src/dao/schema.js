const executor = require("./executor");

module.exports.update = function(data) {
  const query = `
    update schema
    set
      name = $2,
      description = $3,
      search_tokens = to_tsvector($4),
      user_id = $5,
      license = $6
    where id = $1
  `;

  const values = [
    data.id,
    data.name, data.description,
    data.name + " " + data.description,
    data.user_id,
    data.license
  ];

  return executor(query, values);
};

module.exports.create = function(data) {
  const query = `
    insert into
    schema(
      complete, user_id, name, description,
      size_x, size_y, size_z, part_length,
      total_size, total_parts, created, license,
      search_tokens
    )
    values(
      $1, $2, $3, $4,
      $5, $6, $7, $8,
      $9, $10, $11, $12,
      to_tsvector($13)
    )
    returning *
  `;

  const values = [
    false, data.user_id, data.name, data.description || "",
    data.size_x, data.size_y, data.size_z, data.part_length,
    0, 0, Date.now(), data.license || "CC0",
    data.name + " " + data.description
  ];

  return executor(query, values, { single_row: true });
};

module.exports.finalize = function(schema_id) {
  const query = `
    update schema s
    set complete = true,
		total_size = (select sum(length(data)) + sum(length(metadata)) from schemapart sp where sp.schema_id = s.id),
		total_parts = (select count(*) from schemapart sp where sp.schema_id = s.id)
    where id = $1
  `;

  return executor(query, [schema_id], { single_row: true });
};

module.exports.get_by_name = function(name) {
  return executor("select * from schema where name = $1", [name], { single_row: true });
};

module.exports.delete_by_id = function(id) {
  return executor("delete from schema where id = $1", [id], { single_row: true });
};

module.exports.get_by_id = function(id) {
  return executor("select * from schema where id = $1", [id], { single_row: true });
};

module.exports.find_by_user_id = function(user_id) {
  return executor("select * from schema where user_id = $1", [user_id]);
};

module.exports.find_by_user_name = function(user_name) {
  const query = `
    select * from schema
    where user_id = (
      select id
      from public.user
      where name = $1
    )
  `;

  return executor(query, [user_name]);
};

module.exports.find_by_keywords = function(keywords) {
  const query = `
    select *
    from schema
    where search_tokens @@ to_tsquery($1)
    and complete = true
    limit 1000
  `;

  return executor(query, [keywords]);
};


module.exports.find_recent = function(count) {
  const query = `
    select *
    from schema
    where complete = true
		order by created desc
    limit $1
  `;

  return executor(query, [Math.min(count, 250)]);
};

module.exports.get_by_schemaname_and_username = function(schema_name, user_name) {
  const query = `
    select *
    from schema
    where user_id = (select id from public.user where name = $1)
    and name = $2
    limit 1
  `;

  const values = [user_name, schema_name];

  return executor(query, values, { single_row: true });
};
