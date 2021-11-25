import loginStore from '../../store/login.js';
import { get, update } from '../../api/schema.js';

export default {
	props: ["schema", "username"],
	data: function(){
		return {
			edit: false,
			name: this.schema.name,
			// allow edits only if the user-id matches
			can_edit: loginStore.claims && loginStore.claims.user_id == this.schema.user_id
		};
	},
	methods: {
		save: function(){
			if (this.name != this.schema.name){
				get(this.username, this.schema.name)
				.then(schema => {
					schema.name = this.name;
					return update(schema);
				})
				.then(schema => {
					this.schema.name = schema.name;
					this.name = schema.name;
					this.$router.push(`/schema/${loginStore.claims.username}/${schema.name}`);
				});
			}
			this.edit = false;
		},
		abort: function(){
			this.name = this.schema.name;
			this.edit = false;
		}
	},
	template: /*html*/`
		<span>
			<span v-if="!edit">
				{{ schema.name }}
				<a v-if="can_edit" v-on:click="edit=true" class="btn btn-xs btn-primary">
					<i class="fa fa-edit"></i>
				</a>
			</span>
			<span v-if="edit">
				<input type="text"
					class="form-control"
					v-model="name"/>
					<div class="btn-group">
						<a class="btn btn-xs btn-success" v-on:click="save">
							<i class="fa fa-check"></i>
						</a>
						<a v-on:click="abort" class="btn btn-xs btn-danger">
							<i class="fa fa-times"></i>
						</a>
					</div>
			</span>
		</span>
	`
};
