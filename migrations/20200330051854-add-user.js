
exports.up = function(db) {
  return db.createTable('user', {
    id: { type: 'bigint', primaryKey: true, autoIncrement: true },
    name: { type: 'string', notNull: true }
  });
};

exports.down = function(db) {
  return db.dropTable("user");
};
