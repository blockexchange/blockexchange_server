const schema_dao = require("../dao/schema");
const millis_in_one_week = 7*24*60*60*1000;

function job(){
  const now = Date.now();
  schema_dao.delete_old_temp_schemas(now - millis_in_one_week);
}

var handle;

module.exports.start = function(){
  handle = setInterval(job, 60000);
};

module.exports.stop = function(){
  clearInterval(handle);
};
