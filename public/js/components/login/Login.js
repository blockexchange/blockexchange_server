
import Breadcrumb from '../Breadcrumb.js';
import state from './state.js';

const username_input = () => m("input", {
  class: "form-control",
  placeholder: "Username",
  value: state.username,
  oninput: e => state.username = e.target.value
});

const password_input = () => m("input[type=password]", {
  class: "form-control",
  placeholder: "Password",
  value: state.password,
  oninput: e => state.password = e.target.value
});

const login_button = () => m("button", {
  class: "btn btn-primary btn-block",
  disabled: !state.username || !state.password,
  onclick: state.login
}, [
  "Login ",
  m("span", { class: "badge badge-danger" }, state.message)
]);

const logout_button = () => m("button", {
  class: "btn btn-primary btn-block",
  onclick: state.logout
}, "Logout");


export default {
  view: function(){
    return [
			m(Breadcrumb, {
				links: [
          { name: "Home", link: "#!/" },
          { name: "Login", active: true }
        ]
			}),
			m("div", { class: "row"}, [
	      m("div", { class: "col-md-4"}),
	      m("div", { class: "col-md-4"}, [
	        m("form", { class: "" }, [
	          username_input(),
	          password_input(),
	          state.isLoggedIn() ? logout_button() : login_button()
	        ])
	      ]),
	      m("div", { class: "col-md-4"})
    	])
		];
  }
};
