const executor = require("./executor");

module.exports.create = function(name) {
	return executor(`
		insert into
		tag(name)
		values($1)
		returning *
	`, [
		name
	], {
		single_row: true
	});
};

module.exports.find_all = function(){
	return executor(`select * from tag`);
};

module.exports.remove = function(id){
	return executor("delete from tag where id = $1", [id]);
};
