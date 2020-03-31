
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
		description: { type: "text", notNull: true },
    complete: { type: 'boolean', notNull: true },
    size_x: { type: 'int', notNull: true },
    size_y: { type: 'int', notNull: true },
    size_z: { type: 'int', notNull: true },
		created: { type: 'datetime', notNull: true },
    part_length: { type: 'int', notNull: true },
		total_size: { type: 'int', notNull: true },
		total_parts: { type: 'int', notNull: true }
  });
};

exports.down = function(db) {
  return db.dropTable("schema");
};
