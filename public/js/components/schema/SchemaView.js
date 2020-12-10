import { search } from '../../api/searchschema.js';

export default {
	props: ["user_name", "schema_name"],
	data: function(){
		return {
			schema: null
		};
	},
	created: function(){
		console.log("created", this.user_name, this.schema_name);
		search({
			schema_name: this.schema_name,
			user_name: this.user_name
		})
		.then(list => this.schema = list[0]);
	},
	template: /*html*/`
		<div>
			Schema view {{ user_name }} / {{ schema_name }}
			<div v-if="schema">
				Schema-id: {{ schema.id }}
			</div>
		</div>
	`
};
