const executor = require("./executor");

module.exports.get_by_name = function(name) {
  return executor("select * from public.user where name = $1", [name], { single_row: true });
};

module.exports.get_by_id = function(id) {
  return executor("select * from public.user where id = $1", [id], { single_row: true });
};

module.exports.get_all = function() {
  return executor("select * from public.user");
};

module.exports.create = function(data) {
  const query = `
    insert into
    public.user(
      role, name, hash, mail, created
    )
    values(
      $1, $2, $3, $4, $5
    )
    returning *
  `;

  const values = [
    data.role, data.name, data.hash, data.mail, Date.now()
  ];

  return executor(query, values, { single_row: true });
};
