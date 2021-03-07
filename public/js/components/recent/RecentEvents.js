import SchemaCard from './SchemaCard.js';
import { find_recent } from '../../api/searchschema.js';

export default {
	data: function(){
		return {
			changed_schematics: []
		};
	},
	created: function(){
		find_recent(12)
		.then(list => this.changed_schematics = list);
	},
	components: {
		"schema-card": SchemaCard
	},
	template: /*html*/`
		<div>
			<div class="row">
				<div class="col-md-12">
					Recent changes
				</div>
			</div>
			<div class="row">
				<div class="col-md-2" v-for="schema in changed_schematics">
					<schema-card :schema="schema"/>
				</div>
			</div>
		</div>
	`
};
