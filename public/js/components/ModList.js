
const badge = (cl, txt) => m("span", {
	class: `badge badge-${cl}`},
	txt
);

const get_badge_class = mod_name => {
	switch(mod_name) {
		case "ignore":
			return "dange";
		case "default":
			return "success";
		default:
			return "secondary";
	}
};

export default {
  view: function(vnode) {
    const schema = vnode.attrs.schema;
    return m("div", Object.keys(schema.mods).map(mod_name => {
    	return badge(get_badge_class(mod_name), mod_name);
    }));
  }
};
