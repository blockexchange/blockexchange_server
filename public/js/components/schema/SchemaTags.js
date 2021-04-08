import { get_all } from '../../api/tag.js';
import { add } from '../../api/schematag.js';
import Tag from '../tags/Tag.js';

export default {
	props: ["schema"],
	components: {
		"tag-label": Tag
	},
	data: function(){
		return {
			tags: [],
			selected_tag: null
		};
	},
	created: function(){
		get_all().then(t => this.tags = t);
	},
	watch: {
		"selected_tag": function(tag){
			console.log("tag::select", tag);
			add(this.schema.id, tag.id)
			.then(function() {
				console.log("tag added");
			});
		}
	},
	template: /*html*/`
		<div>
			<tag-label v-for="tag in schema.tags" :tag="tag"/>
			<select class="form-control" v-model="selected_tag">
				<option></option>
				<option v-for="tag in tags" :value="tag">
					{{ tag.name }}
				</option>
			</select>
		</div>
	`
};
