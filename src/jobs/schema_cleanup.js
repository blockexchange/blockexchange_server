//const millis_in_one_week = 7*24*60*60*1000;

function job(){
	// TODO: cleanup schemas without description
  //const now = Date.now();
  //const max_age = now - millis_in_one_week;
}

var handle;

module.exports.start = function(){
  handle = setInterval(job, 60000);
};

module.exports.stop = function(){
  clearInterval(handle);
};
