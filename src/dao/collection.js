const executor = require("./executor");

module.exports.create = function(collection) {
	return executor(`
		insert into
		collection(user_id, name)
		values($1, $2)
		returning *
	`, [
		collection.user_id, collection.name
	], {
		single_row: true
	});
};

module.exports.update = function(collection){
	return executor(`
		update collection
		set name = $1
		where id = $2
		and user_id = $3
	`, [
		collection.name,
		collection.id,
		collection.user_id
	], { single_row: true });
};

module.exports.remove = function(user_id, id){
	return executor("delete from collection where user_id = $1 and id = $2", [
		user_id, id
	]);
};

module.exports.find_all_by_userid = function(user_id){
	const query = `
		select * from collection
		where user_id = $1
	`;

	const values = [user_id];
	return executor(query, values);
};

module.exports.find_by_collectionid = function(collection_id){
	const query = `
		select * from collection
		where id = $1
	`;

	const values = [collection_id];
	return executor(query, values, { single_row: true });
};
