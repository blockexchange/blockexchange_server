const executor = require("./executor");

module.exports.create = function(collection_schema) {
	return executor(`
		insert into
		collection_schema(collection_id, schema_id)
		values($1, $2)
		returning *
	`, [
		collection_schema.collection_id,
		collection_schema.schema_id
	], {
		single_row: true
	});
};

module.exports.remove = function(collection_id, schema_id){
	return executor(`
		delete from collection_schema
		where collection_id = $1
		and schema_id = $2
	`, [
		collection_id,
		schema_id
	]);
};

module.exports.find_all_by_collectionid = function(collection_id){
	const query = `
		select * from collection_schema
		where collection_id = $1
	`;

	const values = [collection_id];
	return executor(query, values);
};

module.exports.change_schemaid = function(from_id, to_id){
	return executor(`
		update collection_schema
		set schema_id = $1
		where schema_id = $2
	`, [
			to_id, from_id
	]);
};
