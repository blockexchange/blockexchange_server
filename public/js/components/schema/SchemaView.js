import { search_by_user_and_schemaname } from '../../api/searchschema.js';
import { get_by_schemaid } from '../../api/screenshot.js';
import License from './License.js';

import SchemaDetail from './SchemaDetail.js';
import SchemaMods from './SchemaMods.js';
import SchemaPreview from './SchemaPreview.js';
import SchemaDownload from './SchemaDownload.js';

export default {
	components: {
		"license-badge": License,
		"schema-detail": SchemaDetail,
		"schema-mods": SchemaMods,
		"schema-preview": SchemaPreview,
		"schema-download": SchemaDownload
	},
	props: ["user_name", "schema_name"],
	data: function(){
		return {
			schema: null,
			screenshots: []
		};
	},
	created: function(){
		search_by_user_and_schemaname(this.user_name, this.schema_name)
		.then(schema => {
			this.schema = schema;
			return get_by_schemaid(schema.id);
		})
		.then(screenshots => {
			this.screenshots = screenshots;
		});
	},
	template: /*html*/`
		<div v-if="schema">
			<h3>
			  {{ schema_name }}
			  <small class="text-muted">by {{ user_name }}</small>
				<span v-if="!schema.complete" class="badge badge-danger">
					Incomplete
				</span>
			</h3>
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
							<pre>{{ schema.description }}</pre>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Used nodes</h5>
							<schema-mods :schema="schema"/>
						</div>
					</div>
				</div>
				<div class="col-md-8">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Preview</h5>
							<schema-preview :screenshots="screenshots" :schema="schema"/>
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
