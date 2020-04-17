
export default {
  view: function(vnode) {
    const schema = vnode.attrs.schema;

    return m("div", { class: "card" }, [
      m("div", { class: "card-body" }, [
        m("b", "Usage:"),
        m("p", "/bx_pos1"),
        m("p", `/bx_load ${schema.user.name} ${schema.name}`)
      ])
    ]);
  }
};
