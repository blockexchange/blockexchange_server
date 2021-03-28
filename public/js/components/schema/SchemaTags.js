import { get_all } from '../../api/tag.js';

export default {
	props: ["schema"],
	data: function(){
		return {
			tags: []
		};
	},
	created: function(){
		get_all().then(t => this.tags = t);
	},
	template: /*html*/`
		<div>
			<select class="form-control">
				<option></option>
				<option v-for="tag in tags">
					{{ tag.name }}
				</option>
			</select>
		</div>
	`
};
