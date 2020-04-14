import { get_by_user_and_schemaname } from '../api/searchschema.js';

import SchemaDetail from '../components/SchemaDetail.js';

export default {
  view: function(vnode){
    if (!vnode.state.result && !vnode.state.busy){
      vnode.state.busy = true;
      get_by_user_and_schemaname(vnode.attrs.username, vnode.attrs.schemaname)
      .then(schema => {
        vnode.state.result = schema;
      });

    }

    if (vnode.state.result) {
      return m(SchemaDetail, { schema: vnode.state.result });
    }
  }
};
