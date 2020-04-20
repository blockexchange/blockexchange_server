
import LicenseBadge from '../LicenseBadge.js';
import ModList from '../ModList.js';

import prettybytesize from '../../util/prettybytesize.js';

const entry = (schema, hide_user) => m("tr", {
		class: schema.complete ? "" : "table-danger"
	}, [
	hide_user ? null : m("td", m("a", { href: "#!/schema/" + schema.user.name }, schema.user.name)),
	m("td", m("a", { href: "#!/schema/" + schema.user.name + "/" + schema.name }, schema.name)),
	m("td", [
		moment(+schema.created).format("YYYY-MM-DD HH:mm"),
		" (",
		moment.duration( moment(+schema.created).diff() ).humanize(true),
		")"
	]),
	m("td", schema.downloads),
	m("td", m(LicenseBadge, { license: schema.license })),
	m("td", m("span", { class: "badge badge-secondary" }, prettybytesize(schema.total_size))),
	m("td", schema.size_x + " / " + schema.size_y + " / " + schema.size_z),
	m("td", schema.total_parts),
	m("td", schema.description.substring(0,15)),
	m("td", m(ModList, { schema: schema }) )
]);

const table = (list, hide_user) => m("table", { class: "table table-striped table-condensed" }, [
	m("thead", [
		m("tr", [
			hide_user ? null : m("th", "User"),
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
		list.map(e => entry(e, hide_user))
	])
]);

export default {
	view(vnode) {
		return table(vnode.attrs.list, vnode.attrs.hide_user);
	}
};
