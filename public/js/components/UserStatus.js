import { get_claims } from '../store/token.js';

export default {
  view: function() {
    const claims = get_claims();
    if (!claims){
      return m("div");
    }

    return m("span", { class: "badge badge-light" }, claims.username);
  }
};
