import loginStore from '../../store/login.js';
import { remove } from '../../api/schema.js';


export default {
	props: ["schema"],
	data: function(){
		return {
			ask: false,
			can_edit: loginStore.claims && loginStore.claims.user_id == this.schema.user_id
		};
	},
	methods: {
		do_remove: function(){
			// delete and route to search page
			remove(this.schema.id)
			.then(() => this.$router.push("/search"));
		}
	},
	template: /*html*/`
		<div class="btn-group" v-if="can_edit">
			<button class="btn btn-sm btn-danger" v-if="!ask" v-on:click="ask=true">
				<i class="fa fa-trash"/> Delete
			</button>
			<button class="btn btn-sm btn-danger" v-if="ask" v-on:click="do_remove">
				<i class="fa fa-check"/> Yes
			</button>
			<button class="btn btn-sm btn-success" v-if="ask" v-on:click="ask=false">
				<i class="fa fa-times"/> No
			</button>
		</div>
	`
};
