
export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;

    if (schema.stars > 0){
      // has 1 or more stars
      return [
        schema.stars,
        " ",
        m("i", { class: "fa-star"})
      ];
    } else {
      // no stars
      return m("i", { class: "far fa-star"});
    }
}
};
