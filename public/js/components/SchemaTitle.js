
export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;

    return m("h3", { style: "display: inline;" }, [
    	m("span", { class: "badge badge-primary"}, schema.id),
    	" ",
    	schema.name,
    	" ",
    	m("small", { class: "text-muted" }, "by " + schema.user.name),
      " ",
      schema.complete ? "" : m("span", { class: "badge badge-warning" }, "Incomplete!")
    ]);
  }
};
