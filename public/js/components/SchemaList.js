
const Entry = entry => m("tr", [
	m("td", entry.uid),
	m("td", entry.user_id),
	m("td", new Date(+entry.created).toString()),
	m("td", entry.total_size),
	m("td", entry.size_x + " / " + entry.size_y + " / " + entry.size_z),
	m("td", entry.total_parts),
	m("td", entry.description)
]);

export default list => m("table", { class: "table table-striped table-condensed" }, [
	m("thead", [
		m("tr", [
			m("th", "UID"),
			m("th", "User ID"),
			m("th", "Created"),
			m("th", "Size [bytes]"),
			m("th", "Size [blocks]"),
			m("th", "Parts"),
			m("th", "Description")
		])
	]),
	m("tbody", [
		list.map(Entry)
	])
]);
