import loginStore from '../../store/login.js';
import { get, update } from '../../api/schema.js';

export default {
	props: ["schema", "username"],
	data: function(){
		return {
			edit: false,
			description: this.schema.description,
			// allow edits only if the user-id matches
			can_edit: loginStore.claims && loginStore.claims.user_id == this.schema.user_id
		};
	},
	methods: {
		save: function(){
			get(this.username, this.schema.name)
			.then(schema => {
				schema.description = this.description;
				return update(schema);
			})
			.then(schema => {
				this.schema.description = schema.description;
				this.edit = false;
			});
		},
		abort: function(){
			this.description = this.schema.description;
			this.edit = false;
		}
	},
	template: /*html*/`
		<div>
			<div v-if="edit">
				<textarea class="form-control" v-model="description"></textarea>
				<div class="btn-group">
					<a class="btn btn-xs btn-success" v-on:click="save">
						<i class="fa fa-check"></i>
					</a>
					<a v-on:click="abort" class="btn btn-xs btn-danger">
						<i class="fa fa-times"></i>
					</a>
				</div>
			</div>
			<div v-else>
				<pre>{{ schema.description }}</pre>
				<a v-if="can_edit" v-on:click="edit=true" class="btn btn-xs btn-primary">
					<i class="fa fa-edit"></i>
				</a>
			</div>
		</div>
	`
};
