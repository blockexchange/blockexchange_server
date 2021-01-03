
module.exports = function(permission){
	return function(req, res, next){
		if (!req.claims || !req.claims.permissions){
			// unauthorized
			res.status(403).json({ message: "unauthorized" });
		} else {
			if (req.claims.permissions.indexOf(permission) >= 0){
				next();
			} else {
				res.status(403).json({ message: "permission not found: " + permission });
			}
		}
	};
};
