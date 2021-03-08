
module.exports = function(err, req, res, next){
	console.error(err.stack, typeof(next));
  res.status(500).json({ message: "internal error" });
};
