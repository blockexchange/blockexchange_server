
export default {
	view: function(vnode){
		const links = vnode.attrs.links.map(l => {
			let content = l.name;
			if (l.link){
				content = m("a", { href: l.link }, l.name);
			}

			return m("li", {
				class: "breadcrumb-item" + (l.active ? " active" : ""),
			}, content);
		});
		return m("nav", [
			m("ol", { class: "breadcrumb" }, links)
		]);
	}
};
