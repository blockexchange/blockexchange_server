
export default {
  view(){
    return m("div", { style: "text-align: center;" }, [
			m("img", { src: "pics/blockexchange.png" }),
			m("hr"),
			m("div", [
				m("a", { class: "btn btn-primary", href: "#!/search" }, "Search"),
			]),
			m("div", [
				m("a", { class: "btn btn-secondary", href: "https://github.com/blockexchange" }, [
					m("i", { class: "fab fa-github" }),
					" Source"
				]),
			])
    ]);
  }
};
