
import UserStatus from '../components/UserStatus.js';

const a = (...args) => m("a", ...args);
const div = (...args) => m("div", ...args);
const ul = (...args) => m("ul", ...args);
const li = (...args) => m("li", ...args);

export default {
  view(){
    return m("nav", { class: "navbar fixed-top navbar-expand-lg navbar-dark bg-dark" }, [
      a({class: "navbar-brand", href: "#!/" }, "Block exchange"),
      div({class:"navbar-collapse collapse"}, [
        ul({class: "navbar-nav mr-auto"}, [
          li({class: "nav-item"}, a({class:"nav-link",href:"#!/login"}, "Login")),
          li({class: "nav-item"}, a({class:"nav-link",href:"#!/users"}, "Users")),
          li({class: "nav-item"}, a({class:"nav-link",href:"#!/search"}, "Search"))
        ])
      ]),
      m(UserStatus)
    ]);
  }
};
