
exports.up = function(db) {
  return db.createTable('schema', {
    id: { type: 'bigint', primaryKey: true, autoIncrement: true },
    /*
    user_id: {
      type: "bigint",
      notNull: true,
      foreignKey: {
        name: "schema_user_fk",
        table: "user",
        mapping: "id",
        rules: {
          onDelete: "CASCADE"
        }
      }
    },
    */
    complete: { type: 'boolean', notNull: true },
    size_x: { type: 'int', notNull: true },
    size_y: { type: 'int', notNull: true },
    size_z: { type: 'int', notNull: true },
    part_length: { type: 'int', notNull: true },
  });
};

exports.down = function(db) {
  return db.dropTable("schema");
};
