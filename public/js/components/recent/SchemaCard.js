import { get_by_schemaid } from '../../api/screenshot.js';

export default {
	props: ["schema"],
	data: function(){
		return {
			screenshot: null
		};
	},
	created: function(){
		get_by_schemaid(this.schema.id)
		.then(screenshots => {
			if (screenshots && screenshots.length >= 1){
				this.screenshot = screenshots[0];
			}
		});
	},
	computed: {
		screenshot_url: function(){
			return `api/schema/${this.schema.id}/screenshot/${this.screenshot.id}?height=240&width=360`;
		}
	},
	template: /*html*/`
		<div class="card">
			<img v-if="screenshot" :src="screenshot_url" class="card-img-top img-fluid">
			<div class="card-body">
				<h5 class="card-title">
					<router-link :to="{ name: 'schemapage', params: { schema_name: schema.name, user_name: schema.user.name }}">
						{{ schema.name }}
					</router-link>
					<small class="text-muted">by {{ schema.user.name }}</small>
				</h5>
				<p>
					{{ schema.total_size | prettysize }}; {{ schema.size_x }} / {{ schema.size_y }} / {{ schema.size_z }} nodes
				</p>
				<pre>{{ schema.description }}</pre>
			</div>
		</div>
	`
};
