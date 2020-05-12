
//import { get_claims } from '/../../store/token.js';

export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;
//    const userstars = vnode.attrs.userstars;
//    const claims = get_claims();
//    const has_self_starred = userstars.find(s => s.user_id === claims.user_id);

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
