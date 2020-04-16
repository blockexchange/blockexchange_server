import UserSchemaList from '../components/UserSchemaList.js';

export default {
  view: function(vnode){
    return m(UserSchemaList, { username: vnode.attrs.username });
  }
};
