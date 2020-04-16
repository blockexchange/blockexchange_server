export default {
  view: function(){
    return m("div", { class: "row"}, [
      m("div", { class: "col-md-4"}),
      m("div", { class: "col-md-4"}, [
        m("form", { class: "" }, [
          m("input", { class: "form-control", placeholder: "Username" }),
          m("input[type=password]", { class: "form-control", placeholder: "Password" }),
          m("button", { class: "btn btn-primary btn-block" }, "Login")
        ])
      ]),
      m("div", { class: "col-md-4"})
    ]);
  }
};
