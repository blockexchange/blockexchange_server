
import LicenseBadge from './LicenseBadge.js';
import { get_claims } from '../store/token.js';
import { remove } from '../api/schema.js';

const badge = (cl, txt) => m("span", {
	class: `badge badge-${cl}`},
	txt
);

const modList = mods => m("div", Object.keys(mods).map(mod_name => {
	return badge("secondary", mod_name);
}));

const actions = (schema, i, list) => {
	const claims = get_claims();
	if (claims && claims.user_id == schema.user_id && claims.permissions.schema.delete)
		return m("button", {
			class: "btn btn-sm btn-danger",
			onclick: () => {
				remove(schema)
				.then(() => list.splice(i, 1));
			}
		}, "Delete");
};

const entry = (entry, i, list) => m("tr", [
	m("td", m("a", { href: "#!/schema/" + entry.user.name }, entry.user.name)),
	m("td", m("a", { href: "#!/schema/" + entry.user.name + "/" + entry.name }, entry.name)),
	m("td", [
		moment(+entry.created).format("YYYY-MM-DD HH:mm"),
		" (",
		moment.duration( moment(+entry.created).diff() ).humanize(true),
		")"
	]),
	m("td", entry.downloads),
	m("td", m(LicenseBadge, { license: entry.license })),
	m("td", badge("success", entry.total_size)),
	m("td", entry.size_x + " / " + entry.size_y + " / " + entry.size_z),
	m("td", entry.total_parts),
	m("td", entry.description),
	m("td", modList(entry.mods)),
	m("td", actions(entry, i, list))
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
			m("th", "Mods"),
			m("th", "Actions")
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
