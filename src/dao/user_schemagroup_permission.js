const executor = require("./executor");

module.exports.get_by_user_and_schemagroup_id = function(user_id, schemagroup_id) {
  return executor(
    `
      select * from user_schemagroup_permission
      where user_id = $1 and schemagroup_id = $2
    `,
    [user_id, schemagroup_id],
    { single_row: true }
  );
};

module.exports.get_by_schemagroup_id = function(schemagroup_id) {
  return executor(
    "select * from user_schemagroup_permission where schemagroup_id = $1",
    [schemagroup_id]
  );
};

module.exports.get_by_user_id = function(user_id) {
  return executor(
    "select * from user_schemagroup_permission where user_id = $1",
    [user_id]
  );
};

module.exports.remove = function(data){
  const query = `
  delete from user_schemagroup_permission
  where user_id = $1 and schemagroup_id = $2
  `;
  const values = [
    data.user_id,
    data.schemagroup_id
  ];

  return executor(query, values, { single_row: true});
};

module.exports.get_all = function() {
  return executor("select * from user_schemagroup_permission");
};

module.exports.create_or_update = function(data) {
  const query = `
    insert into
    user_schemagroup_permission(user_id, schemagroup_id, "create", "update", "delete", "manage")
    values($1, $2, $3, $4, $5, $6)
    on conflict(user_id, schemagroup_id)
    do update set
    "create" = EXCLUDED."create",
    "update" = EXCLUDED."update",
    "delete" = EXCLUDED."delete",
    "manage" = EXCLUDED."manage",
    returning *
  `;

  const values = [
    data.user_id,
    data.schemagroup_id,
    data.create,
    data.update,
    data.delete,
    data.manage
  ];

  return executor(query, values, { single_row: true });
};
