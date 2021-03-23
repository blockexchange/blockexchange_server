import loginStore from '../../store/login.js';
import { updateInfo } from '../../api/schema.js';


export default {
	props: ["schema"],
	data: function(){
		return {
			busy: false,
			can_edit: loginStore.claims && loginStore.claims.user_id == this.schema.user_id
		};
	},
	methods: {
		updateInfo: function(){
			this.busy = true;
			updateInfo(this.schema.id)
			.then(() => this.busy = false);
		}
	},
	template: /*html*/`
		<button class="btn btn-sm btn-primary" v-on:click="updateInfo" v-if="can_edit">
			<i v-bind:class="{'fa': true, 'fa-sync': true, 'fa-spin': busy}"/> Update stats and preview
		</button>
	`
};
