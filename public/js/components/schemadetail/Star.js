
export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;

    return [
      m("i", { class: `fa${schema.stars ? "" : "r"} fa-star`})
    ];
}
};
