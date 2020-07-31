import SchemaUsage from '../SchemaUsage.js';
import Star from './Star.js';
import SchemaTitle from '../SchemaTitle.js';
import EditButton from './EditButton.js';
import DeleteButton from './DeleteButton.js';

import html from '../html.js';

import { get_by_user_and_schemaname } from '../../api/searchschema.js';
import { get_all as get_all_stars } from '../../api/userschemastar.js';

import detailtable from './detailtable.js';

export default class {
	constructor(vnode) {
		this.state = {
			progress: 0,
			username: vnode.attrs.username,
			schemaname: vnode.attrs.schemaname,
			schema: null,
			userstars: null,
			ready: false
		};

		this.load_data();
	}

	load_data(){
		get_by_user_and_schemaname(this.state.username, this.state.schemaname)
		.then(s => this.state.schema = s)
		.then(() => get_all_stars(this.state.schema.id))
		.then(userstars => this.state.userstars = userstars)
		.then(() => this.state.ready = true);
	}

  view() {
		const schema = this.state.schema;
		const userstars = this.state.userstars;

		if (!this.state.ready){
			return html`<div>Loading...</div>`;
		}

		return html`
				<div class="row">
					<div class="col-md-8">
						<${SchemaTitle} schema=${schema}/>
						${" "}
						<${Star} schema=${schema} userstars=${userstars}/>
					</div>
					<div class="col-md-4 btn-group" style="text-align: right;">
						<${EditButton} schema=${schema}/>
						<${DeleteButton} schema=${schema}/>
					</div>
				</div>
				<hr/>
				<div class="row">
					<div class="col-md-6">
						<pre>${schema.description || "<no description>"}</pre>
					</div>
					<div class="col-md-6">
						<a class="btn btn-secondary btn-xs" href=${`#!/schema/${this.state.username}/${this.state.schemaname}/preview`}>
							3D Preview <i class="fa fa-play"/>
						</a>
					</div>
				</div>
				<hr/>
				<div class="row">
					<div class="col-md-6">
						${detailtable(schema)}
					</div>
					<div class="col-md-6">
						<${SchemaUsage} schema=${schema}/>
					</div>
				</div>
			`;
  }
}
