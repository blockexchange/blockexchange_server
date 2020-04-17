import LicenseBadge from './LicenseBadge.js';
import Preview from './preview/Preview.js';
import SchemaUsage from './SchemaUsage.js';
import ModList from './ModList.js';

import { get_by_user_and_schemaname } from '../api/searchschema.js';

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
      m("h3", [
        m("span", { class: "badge badge-primary"}, schema.id),
        " ",
        schema.name,
        " ",
        m("small", { class: "text-muted" }, "by " + schema.user.name)
      ]),

			m("div", { class: "row" }, [
				m("div", { class: "col-md-8"}, [
					m("h4", schema.description),
	        schema.long_description
				]),
				m("div", { class: "col-md-4"}, [
					m("div", [
						"Size [bytes]: ", m("span", { class: "badge badge-secondary"}, schema.total_size)
					]),
					"License: ",
					m(LicenseBadge, { license: schema.license }),
					m("div", [
						"Created: ",
						moment(+schema.created).format("YYYY-MM-DD HH:mm"),
						" (",
						moment.duration( moment(+schema.created).diff() ).humanize(true),
						")"
					]),
					m(SchemaUsage, { schema: schema })
				])
			]),

			m("div", { class: "row" }, [
				m("div", { class: "col-md-8" }, [
					m(Preview, { schema: schema, progressCallback: f => this.state.progress = f * 100 }),
					m("div", { class: "progress"}, [
						m("div", { class: "progress-bar", style: `width: ${this.state.progress}%` }, [
							(Math.floor(this.state.progress * 10) / 10) + "%"
						])
					])
				]),
				m("div", { class: "col-md-4" }, [
					"Mod dependencies:",
					m(ModList, { schema: schema })
				])
			])
    ]);
  }
}
