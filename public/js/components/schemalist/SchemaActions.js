import { get_claims } from '../../store/token.js';


export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;
    const removeItem = vnode.attrs.removeItem;

    const claims = get_claims();

    if (claims && claims.user_id == schema.user_id && claims.permissions.schema.delete)
      return m("button", {
        class: "btn btn-sm btn-danger",
        onclick: () => removeItem(schema)
      }, "Delete");
  }
};
