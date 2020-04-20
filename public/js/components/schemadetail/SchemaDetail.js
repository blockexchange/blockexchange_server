import Preview from '../preview/Preview.js';
import SchemaUsage from '../SchemaUsage.js';
import Breadcrumb from '../Breadcrumb.js';
import Star from './Star.js';
import SchemaTitle from '../SchemaTitle.js';
import EditButton from './EditButton.js';

import { get_by_user_and_schemaname } from '../../api/searchschema.js';

import detailtable from './detailtable.js';

export default class {
	constructor(vnode) {
		this.state = {
			progress: 0,
			username: vnode.attrs.username,
			schemaname: vnode.attrs.schemaname,
			schema: null
		};

		this.load_data();
	}

	load_data(){
		get_by_user_and_schemaname(this.state.username, this.state.schemaname)
		.then(s => this.state.schema = s);
	}

  view(vnode) {
    const schema = this.state.schema;

		if (!schema){
			return m("div", "Loading...");
		}

    return m("div", [
			m(Breadcrumb, {
				links: [{
					name: "Home",
					link: "#!/"
				},{
					name: "User-schemas",
				},{
					name: vnode.attrs.username,
					link: "#!/schema/" + vnode.attrs.username
				},{
					name: vnode.attrs.schemaname,
					active: true
				}]
			}),
			m("div", { class: "row" }, [
				m("div", { class: "col-md-8" },	[
					m(SchemaTitle, { schema: schema }),
					" ",
					m(Star, {
						schema: schema,
						load_data: () => this.load_data()
					})
				]),
				m("div", { class: "col-md-4", style: "text-align: right;" }, m(EditButton, {
					schema: schema
				}))
			]),
			m("hr"),
			m("div", { class: "row" }, [
				m("div", { class: "col-md-6"}, [
					m("pre", schema.description || "<no description>")
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
