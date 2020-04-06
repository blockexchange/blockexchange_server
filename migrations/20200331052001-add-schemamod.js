
exports.up = function(db) {
  return db.createTable('schemamod', {
		id: { type: 'bigint', primaryKey: true, autoIncrement: true },
		mod_name: { type: "string", length: 32, notNull: true },
    node_count: { type: "int", notNull: true },
		schema_id: {
      type: "bigint",
      notNull: true,
      foreignKey: {
        name: "schemamod_schema_fk",
        table: "schema",
        mapping: "id",
        rules: {
          onDelete: "CASCADE"
        }
      }
    },
	});
};

exports.down = function(db) {
  return db.dropTable("schemamod");
};
