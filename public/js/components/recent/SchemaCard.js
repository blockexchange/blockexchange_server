import { get_by_schemaid } from '../../api/screenshot.js';
import Tag from '../../components/tags/Tag.js';
import prettysize from '../../util/prettysize.js';

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
	methods: {
		prettysize: prettysize
	},
	computed: {
		screenshot_url: function(){
			return `api/schema/${this.schema.id}/screenshot/${this.screenshot.id}?height=240&width=360`;
		},
		router_link: function(){
			return {
				name: 'schemapage',
				params: {
					schema_name: this.schema.name,
					user_name: this.schema.user.name
				}
			};
		}
	},
	template: /*html*/`
		<div class="card">
			<router-link :to="router_link">
				<img v-if="screenshot" :src="screenshot_url" class="card-img-top img-fluid">
				<div v-if="!screenshot" style="height: 240px; width: 360px"></div>
			</router-link>
			<div class="card-body">
				<h5 class="card-title">
					<p>
						<router-link :to="router_link">
							{{ schema.name }}
						</router-link>
					</p>
					<p>
						<small class="text-muted">by {{ schema.user.name }}</small>
						&nbsp;
						<i v-if="schema.stars > 0" class="fa-regular fa-star"></i>
						<span v-if="schema.stars > 0" class="badge bg-secondary rounded-pill">{{ schema.stars }}</span>
					</p>
				</h5>
				<p>
					<tag-label v-for="tag in schema.tags" :tag_id="tag.tag_id"/>
				</p>
				<p>
					{{ prettysize(schema.total_size) }};
					{{ schema.size_x }} / 
					{{ schema.size_y }} / 
					{{ schema.size_z }} nodes
				</p>
			</div>
		</div>
	`
};
