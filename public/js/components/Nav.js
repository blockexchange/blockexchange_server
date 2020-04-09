

const a = (...args) => m("a", ...args);
const div = (...args) => m("div", ...args);
const form = (...args) => m("form", ...args);
const input = (...args) => m("input", ...args);
const button = (...args) => m("button", ...args);
const ul = (...args) => m("ul", ...args);
const li = (...args) => m("li", ...args);

export default {
  view(){
    return m("nav", { class: "navbar fixed-top navbar-light bg-light" }, [
      a({class: "navbar-brand", href: "#!/" }, "Block exchange"),
      div([
        ul({class: "navbar-nav"}, [
          li({class: "nav-item"}, a({class:"nav-link",href:"#!/search"}, "Search"))
        ])
      ]),
      form({class: "form-inline"}, [
        input({class: "form-control mr-sm-2", placeholder: "Search"}),
        button({class: "btn btn-outline-success my-2 my-sm-0"}, "Search")
      ])
    ]);
  }
};
