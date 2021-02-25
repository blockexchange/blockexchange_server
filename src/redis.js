const redis = require("redis");

let client = null;

if (process.env.REDIS_HOST){
	// redis available, create a client
	client = redis.createClient({
		host: process.env.REDIS_HOST,
		port: +process.env.REDIS_PORT,
		detect_buffers: true
	});
}

module.exports = client;
