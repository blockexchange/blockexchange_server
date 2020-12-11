import { search_by_user_and_schemaname } from '../../api/searchschema.js';

export default {
	props: ["user_name", "schema_name"],
	data: function(){
		return {
			schema: null
		};
	},
	created: function(){
		search_by_user_and_schemaname(this.user_name, this.schema_name)
		.then(s => this.schema = s);
	},
	template: /*html*/`
		<div>
			<h3>
			  {{ schema_name }}
			  <small class="text-muted">by {{ user_name }}</small>
			</h3>
			<div class="row" v-if="schema">
				<div class="col-md-12">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Images</h5>
							...
						</div>
					</div>
				</div>
			</div>
			<br>
			<div class="row" v-if="schema">
				<div class="col-md-8">
					<div class="card">
					  <div class="card-body">
					    <h5 class="card-title">Description</h5>
					    <pre>{{ schema.description }}</pre>
					  </div>
					</div>
				</div>
				<div class="col-md-4">
					<div class="card">
						<div class="card-body">
							<h5 class="card-title">Details</h5>
							<ul>
								<li>
									<b>Changed: </b>{{ new Date(+schema.created).toDateString() }}
								</li>
								<li>
									<b>Size: </b>{{ schema.total_size }} bytes
								</li>
								<li>
									<b>Dimensions: </b>{{ schema.max_x }} / {{ schema.max_y }} / {{ schema.max_z }} blocks
								</li>
								<li>
									<b>Parts: </b>{{ schema.total_parts }}
								</li>
							</ul>
						</div>
					</div>
				</div>
			</div>
			<br>
			<div class="row" v-if="schema">
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
