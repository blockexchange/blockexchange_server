
exports.up = function(db) {
  return db.createTable('schemapart', {
    id: { type: 'bigint', primaryKey: true, autoIncrement: true },
    schema_id: {
      type: "bigint",
      notNull: true,
      foreignKey: {
        name: "schema_schemapart_fk",
        table: "schema",
        mapping: "id",
        rules: {
          onDelete: "CASCADE"
        }
      }
    },
    offset_x: { type: 'int', notNull: true },
    offset_y: { type: 'int', notNull: true },
    offset_z: { type: 'int', notNull: true },
    data: { type: 'blob', notNull: true },
  });
};

exports.down = function(db) {
  return db.dropTable("schemapart");
};
