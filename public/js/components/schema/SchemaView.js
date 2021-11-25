import { search_by_user_and_schemaname } from '../../api/searchschema.js';
import { get_by_schemaid } from '../../api/screenshot.js';

import loginStore from '../../store/login.js';

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
		"schema-star": SchemaStar
	},
	props: ["user_name", "schema_name"],
	data: function(){
		return {
			schema: null,
			screenshots: [],
			is_owner: false,
			preview_version: 0
		};
	},
	methods: {
		update: function(){
			search_by_user_and_schemaname(this.user_name, this.schema_name)
			.then(schema => {
				this.schema = schema;
				this.is_owner = loginStore.claims && loginStore.claims.user_id == schema.user_id;
				return get_by_schemaid(schema.id);
			})
			.then(screenshots => {
				this.screenshots = screenshots;
			});
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
						<schema-title :schema="schema" :username="user_name"/>
						<small class="text-muted">by {{ user_name }}</small>
						<span v-if="!schema.complete" class="badge bg-danger">
							Incomplete
						</span>
					</h3>
					<schema-star :schema="schema"/>
				</div>
				<div class="col-md-6">
					<div class="btn-group float-end">
						<schema-updateinfo :schema="schema" :username="user_name" v-on:stats-updated="preview_version++"/>
						<schema-delete :schema="schema" :username="user_name"/>
					</div>
				</div>
			</div>
			<div class="row">
				<div class="col-md-4">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Details</h5>
							<schema-detail :schema="schema" :username="user_name"/>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Description</h5>
							<schema-description :schema="schema" :username="user_name"/>
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
