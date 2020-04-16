import SchemaList from '../components/SchemaList.js';

import { find_by_username } from '../api/searchschema.js';

export default {
  oninit: function(vnode) {
    vnode.state.list = [];
    find_by_username(vnode.attrs.username)
    .then(l => vnode.state.list = l);
  },
  view: function(vnode){
    return m(SchemaList, { list: vnode.state.list });
  }
};
