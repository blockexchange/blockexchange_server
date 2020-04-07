
exports.up = function(db) {
  return db.createTable('schema', {
    id: { type: 'bigint', primaryKey: true, autoIncrement: true },
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
    name: { type: "string", length: 64, notNull: true },
		description: { type: "text", notNull: true },
    complete: { type: 'boolean', notNull: true },
    size_x: { type: 'smallint', notNull: true },
    size_y: { type: 'smallint', notNull: true },
    size_z: { type: 'smallint', notNull: true },
		created: { type: 'bigint', notNull: true },
    part_length: { type: 'smallint', notNull: true },
		total_size: { type: 'int', notNull: true },
		total_parts: { type: 'int', notNull: true },
    search_tokens: { type: "TSVECTOR", notNull: true }
  })
  .then(() => {
    return db.addIndex("schema", "schema_uid", ["uid"], true)
		.then(() => db.addIndex("schema", "schema_created", ["created"]));
  });
};

exports.down = function(db) {
  return db.dropTable("schema");
};
