
export default {
  view: function(vnode){
    return JSON.stringify(vnode.attrs.schema);
  }
};
