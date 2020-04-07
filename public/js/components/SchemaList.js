
const ModList = mods => m("div", Object.keys(mods).map(mod_name => {
	return m("span", { class: "badge badge-secondary" }, mod_name);
}));

const Entry = entry => m("tr", [
	m("td", entry.user.name),
	m("td", entry.name),
	m("td", new Date(+entry.created).toString()),
	m("td", entry.total_size),
	m("td", entry.size_x + " / " + entry.size_y + " / " + entry.size_z),
	m("td", entry.total_parts),
	m("td", entry.description),
	m("td", ModList(entry.mods))
]);

export default list => m("table", { class: "table table-striped table-condensed" }, [
	m("thead", [
		m("tr", [
			m("th", "User"),
			m("th", "Name"),
			m("th", "Created"),
			m("th", "Size [bytes]"),
			m("th", "Size [blocks]"),
			m("th", "Parts"),
			m("th", "Description"),
			m("th", "Mods")
		])
	]),
	m("tbody", [
		list.map(Entry)
	])
]);
