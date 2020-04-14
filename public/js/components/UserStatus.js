const form = (...args) => m("form", ...args);
const input = (...args) => m("input", ...args);
const input_password = (...args) => m("input[type=password]", ...args);
const button = (...args) => m("button", ...args);

export default {
  view: function() {
    return form({class: "form-inline"}, [
      input({class: "form-control mr-sm-2", placeholder: "Username"}),
      input_password({class: "form-control mr-sm-2", placeholder: "Password"}),
      button({class: "btn btn-outline-success my-2 my-sm-0"}, "Login")
    ]);
  }
};
