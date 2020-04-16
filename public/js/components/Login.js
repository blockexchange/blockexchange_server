
import { get_claims, set_token } from '../store/token.js';
import { request_token } from '../api/token.js';

const store = {
  username: "",
  password: ""
};

const username_input = () => m("input", {
  class: "form-control",
  placeholder: "Username",
  value: store.username,
  oninput: e => store.username = e.target.value
});

const password_input = () => m("input[type=password]", {
  class: "form-control",
  placeholder: "Password",
  value: store.password,
  oninput: e => store.password = e.target.value
});

const login_button = () => m("button", {
  class: "btn btn-primary btn-block",
  disabled: !store.username || !store.password,
  onclick: () => {
    request_token(store.username, store.password)
    .then(token => set_token(token))
    .catch(e => console.log(e));
  }
}, "Login");

const logout_button = () => m("button", {
  class: "btn btn-primary btn-block",
  onclick: () => set_token(null)
}, "Logout");



export default {
  view: function(){
    return m("div", { class: "row"}, [
      m("div", { class: "col-md-4"}),
      m("div", { class: "col-md-4"}, [
        m("form", { class: "" }, [
          username_input(),
          password_input(),
          get_claims() ? logout_button() : login_button()
        ])
      ]),
      m("div", { class: "col-md-4"})
    ]);
  }
};
