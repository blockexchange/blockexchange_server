const executor = require("./executor");

module.exports.get_by_name = function(name) {
  return executor("select * from public.user where name = $1", [name], { single_row: true });
};

module.exports.get_by_id = function(id) {
  return executor("select * from public.user where id = $1", [id], { single_row: true });
};

module.exports.get_by_external_id = function(id) {
  return executor("select * from public.user where external_id = $1", [id], { single_row: true });
};

module.exports.get_all = function() {
  return executor("select * from public.user");
};

module.exports.create = function(data) {
  const query = `
    insert into
    public.user(
      role, name, hash, mail, type, external_id, created
    )
    values(
      $1, $2, $3, $4, $5, $6, $7
    )
    returning *
  `;

  const values = [
    data.role, data.name, data.hash, data.mail, data.type, data.external_id, Date.now()
  ];

  return executor(query, values, { single_row: true });
};

module.exports.update_user = function(user){
	return executor(`
		update public.user
		set name = $2, mail = $3
		where id = $1
	`, [
		user.id, user.name, user.mail
	]);
};
