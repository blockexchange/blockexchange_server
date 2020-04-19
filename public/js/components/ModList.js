
const badge = (cl, txt) => m("span", {
	class: `badge badge-${cl}`},
	txt
);

export default {
  view: function(vnode) {
    const schema = vnode.attrs.schema;
    return m("div", Object.keys(schema.mods).map(mod_name => {
    	return badge(mod_name == "ignore" ? "danger" : "secondary", mod_name);
    }));
  }
};
