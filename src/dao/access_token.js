const executor = require("./executor");

module.exports.create = function(user_id, created, expires, name, token) {
  const query = `
    insert into
    access_token(user_id, created, expires, name, token)
    values($1, $2, $3, $4, $5)
    returning *
  `;

  const values = [
    user_id, created, expires, name, token
  ];

  return executor(query, values, { single_row: true });
};

module.exports.remove = function(user_id, id){
	return executor("delete from access_token where user_id = $1 and id = $2", [user_id, id]);
};

module.exports.find_all_by_userid = function(user_id) {
  const query = `
    select * from access_token
    where user_id = $1
  `;

  const values = [
    user_id
  ];

  return executor(query, values);
};

module.exports.increment_usecount = function(id) {
  const query = `
    update access_token
    set usecount = usecount + 1
    where id = $1
  `;

  const values = [id];

  return executor(query, values, { single_row: true });
};
