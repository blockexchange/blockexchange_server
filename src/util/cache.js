const redis = require("../redis");

module.exports.set = function(key, value){
	if (redis){
		// populate cache
		redis.set(key, value);
		// let data expire after a while
		redis.set(key, value, 'EX', 60 * 60 * 24);
	}
};

module.exports.get = function(key){
	if (redis){
		// use cache
		return new Promise((resolve, reject) => {
			redis.get(key, function(err, reply) {
				if (err){
					reject(err);
				} else {
					resolve(reply);
				}
			});
		});
	} else {
		// no cache
		return Promise.resolve(null);
	}
};
