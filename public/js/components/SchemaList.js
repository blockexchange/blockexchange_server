
const licenseBadge = license => {
	switch (license) {
		case "CC0":
			return m("img", {src:"pics/license_cc0.png"});
		default:
			return license;
	}
};

const badge = (cl, txt) => m("span", {
	class: `badge badge-${cl}`},
	txt
);

const modList = mods => m("div", Object.keys(mods).map(mod_name => {
	return badge("secondary", mod_name);
}));

const entry = entry => m("tr", [
	m("td", m("a", { href: "#!/schema/" + entry.user.name }, entry.user.name)),
	m("td", m("a", { href: "#!/schema/" + entry.user.name + "/" + entry.name }, entry.name)),
	m("td", [
		moment(+entry.created).format("YYYY-MM-DD HH:mm"),
		" (",
		moment.duration( moment(+entry.created).diff() ).humanize(true),
		")"
	]),
	m("td", entry.downloads),
	m("td", licenseBadge(entry.license)),
	m("td", badge("success", entry.total_size)),
	m("td", entry.size_x + " / " + entry.size_y + " / " + entry.size_z),
	m("td", entry.total_parts),
	m("td", entry.description),
	m("td", modList(entry.mods))
]);

const table = list => m("table", { class: "table table-striped table-condensed" }, [
	m("thead", [
		m("tr", [
			m("th", "User"),
			m("th", "Name"),
			m("th", "Created"),
			m("th", "Downloads"),
			m("th", "License"),
			m("th", "Size [bytes]"),
			m("th", "Size [blocks]"),
			m("th", "Parts"),
			m("th", "Description"),
			m("th", "Mods")
		])
	]),
	m("tbody", [
		list.map(entry)
	])
]);

export default {
	view(vnode) {
		return table(vnode.attrs.list);
	}
};
