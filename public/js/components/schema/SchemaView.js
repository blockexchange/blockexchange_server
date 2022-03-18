import { search } from '../../api/searchschema.js';
import { get_by_schemaid } from '../../api/screenshot.js';

import loginStore from '../../store/login.js';
import infoStore from '../../store/info.js';

import SchemaDetail from './SchemaDetail.js';
import SchemaMods from './SchemaMods.js';
import SchemaPreview from './SchemaPreview.js';
import SchemaDownload from './SchemaDownload.js';
import SchemaDescription from './SchemaDescription.js';
import SchemaDelete from './SchemaDelete.js';
import SchemaUpdateInfo from './SchemaUpdateInfo.js';
import SchemaTitle from './SchemaTitle.js';
import SchemaTags from './SchemaTags.js';
import SchemaStar from './SchemaStar.js';

import Clipboard from '../Clipboard.js';

export default {
	components: {
		"schema-detail": SchemaDetail,
		"schema-mods": SchemaMods,
		"schema-preview": SchemaPreview,
		"schema-download": SchemaDownload,
		"schema-description": SchemaDescription,
		"schema-delete": SchemaDelete,
		"schema-title": SchemaTitle,
		"schema-tags": SchemaTags,
		"schema-updateinfo": SchemaUpdateInfo,
		"schema-star": SchemaStar,
		"clipboard-component": Clipboard
	},
	props: ["user_name", "schema_name"],
	data: function(){
		return {
			schema: null,
			screenshots: [],
			is_owner: false,
			preview_version: Date.now()
		};
	},
	methods: {
		update: function(){
			search({
				user_name: this.user_name,
				schema_name: this.schema_name
			})
			.then(result => result.list[0])
			.then(schema => {
				this.schema = schema;
				this.is_owner = loginStore.claims && loginStore.claims.user_id == schema.user_id;
				return get_by_schemaid(schema.id);
			})
			.then(screenshots => {
				this.screenshots = screenshots;
			});
		},
		updatePreview: function(){
			this.preview_version = Date.now();
		},
		getLink: function(){
			//TODO: move base_url to common props
			return `${infoStore.oauth.base_url}/api/static/schema/${this.user_name}/${this.schema_name}`;
		}
	},
	created: function(){
		this.update();
	},
	template: /*html*/`
		<div v-if="schema">
			<div class="row">
				<div class="col-md-6">
					<h3>
						<schema-title :schema="schema"/>
						<small class="text-muted">by {{ user_name }}</small>
						<span v-if="!schema.complete" class="badge bg-danger">
							Incomplete
						</span>
						&nbsp;
						<clipboard-component :link="getLink()"/>
						&nbsp;
						<schema-star :schema="schema"/>
						</h3>
				</div>
				<div class="col-md-6">
					<div class="btn-group float-end">
						<schema-updateinfo :schema="schema" v-on:stats-updated="updatePreview()"/>
						<schema-delete :schema="schema"/>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-md-4">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Details</h5>
							<schema-detail :schema="schema"/>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Description</h5>
							<schema-description :schema="schema"/>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Used mods</h5>
							<schema-mods :schema="schema"/>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Tags</h5>
							<schema-tags :schema="schema" :is_owner="is_owner" v-on:updated="update"/>
						</div>
					</div>
				</div>
				<div class="col-md-8">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Preview</h5>
							<schema-preview :screenshots="screenshots" :schema="schema" :version="preview_version"/>
						</div>
					</div>
				</div>
			</div>
			<br>
			<div class="row">
				<div class="col-md-12">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Download</h5>
							<schema-download :schema="schema"/>
						</div>
					</div>
				</div>
			</div>
		</div>
	`
};
