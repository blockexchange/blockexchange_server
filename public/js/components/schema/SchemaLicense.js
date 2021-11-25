import loginStore from '../../store/login.js';
import { get, update } from '../../api/schema.js';

export default {
	props: ["schema", "username"],
	data: function(){
		return {
			edit: false,
			license: this.schema.license,
			// allow edits only if the user-id matches
			can_edit: loginStore.claims && loginStore.claims.user_id == this.schema.user_id
		};
	},
	methods: {
		save: function(){
			get(this.username, this.schema.name)
			.then(schema => {
				schema.license = this.license;
				return update(schema);
			})
			.then(schema => {
				this.schema.license = schema.license;
				this.edit = false;
			});
		},
		abort: function(){
			this.license = this.schema.license;
			this.edit = false;
		}
	},
	computed: {
		imgsrc: function(){
			switch (this.license){
				case "CC0": return "pics/license_cc0.png";
				case "CC-BY-SA 3.0": return "pics/license_cc-by-sa.png";
			}
		}
	},
	template: /*html*/`
		<div>
			<div v-if="edit">
				<input type="text" class="form-control" v-model="license"/>
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
				<img v-if="imgsrc" :src="imgsrc"/>
				<div v-else class="badge bg-primary">
					{{ license }}
				</div>
				<a v-if="can_edit" v-on:click="edit=true" class="btn btn-xs btn-primary">
					<i class="fa fa-edit"></i>
				</a>
			</div>
		</div>
	`
};
