
function job(){

}

var handle;

module.exports.start = function(){
  handle = setInterval(job, 60000);
};

module.exports.stop = function(){
  clearInterval(handle);
};
