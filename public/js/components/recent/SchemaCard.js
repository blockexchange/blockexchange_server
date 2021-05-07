import { get_by_schemaid } from '../../api/screenshot.js';
import Tag from '../../components/tags/Tag.js';

export default {
	props: ["schema"],
	data: function(){
		return {
			screenshot: null
		};
	},
	components: {
		"tag-label": Tag
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
					<tag-label v-for="tag in schema.tags" :tag_id="tag.tag_id"/>
				</p>
				<p>
					{{ schema.total_size | prettysize }};
					{{ schema.size_x_plus+schema.size_x_minus }} / 
					{{ schema.size_y_plus+schema.size_y_minus }} / 
					{{ schema.size_z_plus+schema.size_z_minus }} nodes
				</p>
				<pre>{{ schema.description }}</pre>
			</div>
		</div>
	`
};
