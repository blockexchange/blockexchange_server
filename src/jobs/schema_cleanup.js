const millis_in_one_week = 7*24*60*60*1000;
const { delete_incomplete_schemas } = require("../dao/schema");

function job(){
	// cleanup old, incomplete schemas
  const now = Date.now();
  const max_age = now - millis_in_one_week;

	delete_incomplete_schemas(max_age);
}

var handle;

module.exports.start = function(){
  handle = setInterval(job, 60000);
};

module.exports.stop = function(){
  clearInterval(handle);
};
