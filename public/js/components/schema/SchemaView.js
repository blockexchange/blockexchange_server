import { search_by_user_and_schemaname } from '../../api/searchschema.js';
import { get_by_schemaid } from '../../api/screenshot.js';
import License from './License.js';

export default {
	components: {
		"license-badge": License
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
							<ul>
								<li>
									<b>Changed: </b>{{ new Date(+schema.created).toDateString() }}
								</li>
								<li>
									<b>Size: </b>{{ schema.total_size | prettysize }}
								</li>
								<li>
									<b>Dimensions: </b>{{ schema.max_x }} / {{ schema.max_y }} / {{ schema.max_z }} blocks
								</li>
								<li>
									<b>Parts: </b>{{ schema.total_parts }}
								</li>
								<li>
									<b>Downloads: </b>{{ schema.downloads }}
								</li>
								<li>
									<b>License: </b><license-badge style="display: inline;" :license="schema.license"/>
								</li>
							</ul>
						</div>
					</div>
					<br>
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Description</h5>
							<pre>{{ schema.description }}</pre>
						</div>
					</div>
				</div>
				<div class="col-md-8">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Preview</h5>
							<div v-for="screenshot in screenshots">
								<img class="img-fluid" :src="'api/schema/' + schema.id + '/screenshot/' + screenshot.id"/>
							</div>
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
							<pre>Download steps here</pre>
						</div>
					</div>
				</div>
			</div>
		</div>
	`
};
