import Preview from '../preview/Preview.js';
import SchemaUsage from '../SchemaUsage.js';

import { get_by_user_and_schemaname } from '../../api/searchschema.js';

import detailtable from './detailtable.js';

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
	        m("pre", schema.long_description || "<no description>")
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
				m("div", { class: "col-md-6"}, detailtable(schema)),
				m("div", { class: "col-md-6"}, [
					m(SchemaUsage, { schema: schema })
				]),
			])

    ]);
  }
}
