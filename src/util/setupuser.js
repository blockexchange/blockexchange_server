const { create, generateToken } = require("../dao/access_token");

// Sets up the existing user with some defaults (access_token, etc)
module.exports = async function(user){
	// generate a "default" access_token with half a year expiration
	await create(user.id, Date.now(), Date.now() + (1000*3600*24*31*6), "default", generateToken());
	return user;
};
