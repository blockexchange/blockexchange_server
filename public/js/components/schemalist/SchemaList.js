
import LicenseBadge from '../LicenseBadge.js';
import ModList from '../ModList.js';
import SchemaActions from './SchemaActions.js';

import prettybytesize from '../../util/prettybytesize.js';


const entry = (entry, removeItem, hide_user) => m("tr", [
	hide_user ? null : m("td", m("a", { href: "#!/schema/" + entry.user.name }, entry.user.name)),
	m("td", m("a", { href: "#!/schema/" + entry.user.name + "/" + entry.name }, entry.name)),
	m("td", [
		moment(+entry.created).format("YYYY-MM-DD HH:mm"),
		" (",
		moment.duration( moment(+entry.created).diff() ).humanize(true),
		")"
	]),
	m("td", entry.downloads),
	m("td", m(LicenseBadge, { license: entry.license })),
	m("td", m("span", { class: "badge badge-secondary" }, prettybytesize(entry.total_size))),
	m("td", entry.size_x + " / " + entry.size_y + " / " + entry.size_z),
	m("td", entry.total_parts),
	m("td", entry.description.substring(0,15)),
	m("td", m(ModList, { schema: entry }) ),
	m("td", m(SchemaActions, { schema: entry, removeItem: removeItem }))
]);

const table = (list, removeItem, hide_user) => m("table", { class: "table table-striped table-condensed" }, [
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
			m("th", "Mods"),
			m("th", "Actions")
		])
	]),
	m("tbody", [
		list.map(e => entry(e, removeItem, hide_user))
	])
]);

export default {
	view(vnode) {
		return table(vnode.attrs.list, vnode.attrs.removeItem, vnode.attrs.hide_user);
	}
};
