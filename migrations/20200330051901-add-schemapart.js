
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
    offset_x: { type: 'smallint', notNull: true },
    offset_y: { type: 'smallint', notNull: true },
    offset_z: { type: 'smallint', notNull: true },
    data: { type: 'blob', notNull: true },
    metadata: { type: 'blob', notNull: true }
  })
  .then(() => db.addIndex("schemapart", "schemapart_coords", [
    "schema_id", "offset_x", "offset_y", "offset_z"
  ]));
};

exports.down = function(db) {
  return db.dropTable("schemapart");
};
