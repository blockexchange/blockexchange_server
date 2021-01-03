const { create, generateToken } = require("../dao/access_token");

// Sets up the existing user with some defaults (access_token, etc)
module.exports = async function(user){
	// generate a "default" access_token
	await create(user.id, Date.now(), Date.now() + (1000*3600*24*7), "default", generateToken());
	return user;
};
