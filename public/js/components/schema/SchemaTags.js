import TagStore from '../../store/tag.js';

import { add } from '../../api/schematag.js';
import Tag from '../tags/Tag.js';

export default {
	props: ["schema", "is_owner"],
	components: {
		"tag-label": Tag
	},
	data: function(){
		return {
			tags: TagStore.tags,
			selected_tag: null
		};
	},
	methods: {
		tags_updated: function(){
			this.$emit("updated");
		}
	},
	watch: {
		"selected_tag": function(tag){
			add(this.schema.id, tag.id)
			.then(() => this.tags_updated());
		}
	},
	template: /*html*/`
		<div>
			<tag-label v-for="tag in schema.tags"
				:tag_id="tag.tag_id"
				:schema_id="schema.id"
				:user_id="schema.user_id"
				v-on:removed="tags_updated"
			/>
			<select class="form-control" v-if="is_owner" v-model="selected_tag">
				<option></option>
				<option v-for="tag in tags" :value="tag">
					{{ tag.name }}
				</option>
			</select>
		</div>
	`
};
