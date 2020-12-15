
module.exports.can_upload = function(rolename){
	return rolename === "ADMIN" || rolename === "UPLOAD_ONLY" || rolename === "MEMBER";
};

module.exports.can_edit = function(rolename){
	return rolename === "ADMIN" || rolename === "MEMBER";
};

module.exports.can_delete = function(rolename){
	return rolename === "ADMIN" || rolename === "MEMBER";
};
