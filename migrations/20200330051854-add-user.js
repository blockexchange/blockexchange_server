
exports.up = function(db) {
  return db.createTable('user', {
    id: { type: 'bigint', primaryKey: true, autoIncrement: true },
    name: { type: 'string', notNull: true },
    hash: { type: 'string', notNull: true },
    mail: { type: 'string' }
  })
  .then(() => {
    return db.addIndex("user", "user_name", ["name"], true);
  });
};

exports.down = function(db) {
  return db.dropTable("user");
};
