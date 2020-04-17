import LicenseBadge from './LicenseBadge.js';
import Preview from './preview/Preview.js';
import SchemaUsage from './SchemaUsage.js';
import ModList from './ModList.js';

import { get_by_user_and_schemaname } from '../api/searchschema.js';

const modDetails = schema => m("ul",
	m("div", Object.keys(schema.mods).map(
		mod_name => m("li", `${mod_name}: ${schema.mods[mod_name]}`)
	)));

const title = schema => m("h3", [
	m("span", { class: "badge badge-primary"}, schema.id),
	" ",
	schema.name,
	" ",
	m("small", { class: "text-muted" }, "by " + schema.user.name)
]);

export default class {
	constructor(vnode) {
		this.state = {
			progress: 0
		};

		get_by_user_and_schemaname(vnode.attrs.username, vnode.attrs.schemaname)
		.then(s => this.state.schema = s);
	}


  view() {
		// TODO: cleanup/separate
    const schema = this.state.schema;

		if (!schema){
			return m("div");
		}

    return m("div", [
			title(schema),
			m("hr"),
			m("div", { class: "row" }, [
				m("div", { class: "col-md-6"}, [
					m("h4", schema.description),
	        schema.long_description
				]),
				m("div", { class: "col-md-6"}, [
					m("div", [
						m(Preview, { schema: schema, progressCallback: f => this.state.progress = f * 100 }),
						m("div", { class: "progress"}, [
							m("div", { class: "progress-bar", style: `width: ${this.state.progress}%` }, [
								(Math.floor(this.state.progress * 10) / 10) + "%"
							])
						])
					])
				])
			]),
			m("hr"),

			m("div", { class: "row" }, [
				m("div", { class: "col-md-6"}, [
					m("table", { class: "table table-condensed table-striped" }, [
						m("tbody", [
							m("tr", [
								m("td", "Size [bytes]"),
								m("td", m("span", { class: "badge badge-secondary"}, schema.total_size))
							]),
							m("tr", [
								m("td", "Size [blocks]"),
								m("td", `${schema.size_x} / ${schema.size_y} / ${schema.size_z}`)
							]),
							m("tr", [
								m("td", "Volume [blocks]"),
								m("td", `${schema.size_x * schema.size_y * schema.size_z}`)
							]),
							m("tr", [
								m("td", "License"),
								m("td", m(LicenseBadge, { license: schema.license }))
							]),
							m("tr", [
								m("td", "Created"),
								m("td", [
									moment(+schema.created).format("YYYY-MM-DD HH:mm"),
									" (",
									moment.duration( moment(+schema.created).diff() ).humanize(true),
									")"
								])
							]),
							m("tr", [
								m("td", "Parts"),
								m("td", schema.total_parts)
							]),
							m("tr", [
								m("td", "Mods"),
								m("td", m(ModList, { schema: schema }))
							]),
							m("tr", [
								m("td", "Mod block count"),
								m("td", modDetails(schema))
							])
						])
					])
				]),
				m("div", { class: "col-md-6"}, [
					m(SchemaUsage, { schema: schema })
				]),
			])

    ]);
  }
}
