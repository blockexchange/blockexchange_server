
exports.up = function(db) {
  return db.createTable('schematag', {
		id: { type: 'bigint', primaryKey: true, autoIncrement: true },
		tag_name: { type: "string", notNull: true },
		schema_id: {
      type: "bigint",
      notNull: true,
      foreignKey: {
        name: "schematag_schema_fk",
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
  return db.dropTable("schematag");
};
