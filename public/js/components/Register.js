import Breadcrumb from './Breadcrumb.js';


export default {
  view: function(){
    return [
			m(Breadcrumb, {
				links: [{
					name: "Home",
					link: "#!/"
				},{
					name: "Register",
					active: true
				}]
			}),
			m("div", { class: "row"}, [
	      m("div", { class: "col-md-4"}),
	      m("div", { class: "col-md-4"}, [
	        m("form", { class: "" }, [
	          m("input", { class: "form-control", placeholder: "Username" }),
	          m("input[type=password]", { class: "form-control", placeholder: "Password" }),
	          m("input[type=password]", { class: "form-control", placeholder: "Password" }),
	          m("button", { class: "btn btn-primary btn-block" }, "Register")
	        ])
	      ]),
	      m("div", { class: "col-md-4"})
	    ])
		];
  }
};
