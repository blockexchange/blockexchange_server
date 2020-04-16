import { get_all } from '../api/user.js';

var list;

export default {
  oninit: function() {
    if (!list){
      get_all()
      .then(l => list = l);
    }
  },
  view: function(){
    if (!list){
      return m("div");
    }
    return m("ul", list.map(user => {
      return m("li", [
        m("a", { href: `#!/schema/${user.name}` }, user.name)
      ]);
    }));
  }
};
