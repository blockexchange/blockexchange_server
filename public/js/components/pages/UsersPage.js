import { get } from '../../api/user.js';
import Pager from '../Pager.js';

export default {
	created: function(){
		this.update();
	},
	components: {
		"pager-component": Pager
	},
	data: function(){
		return {
			userdata: null,
			page: 1,
			limit: 20
		};
	},
	watch: {
		"$route": "update"
	},
	methods: {
		update: function(){
			const r = this.$route;
			this.page = r.query.page || 1;
			get(this.limit, (this.page-1) * this.limit)
			.then(userdata => this.userdata = userdata);
		},
		updateRoute: function(page){
			const r = this.$router.currentRoute.value;
			this.$router.push({
				path: r.path,
				query: {
					page: page
				}
			});
		}
	},
	template: /*html*/`
		<div v-if="userdata">
			<pager-component
				:current="this.page"
				:pages="Math.ceil(this.userdata.total / this.limit)"
				v-on:switch="updateRoute">
			</pager-component>
			<table class="table table-condensed table-striped">
				<thead>
					<tr>
						<th>Name</th>
						<th>Created</th>
						<th>Type</th>
						<th>Role</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="user in userdata.list">
						<td>{{ user.name }}</td>
						<td>{{ new Date(+user.created).toDateString() }}</td>
						<td>{{ user.type }}</td>
						<td>{{ user.role }}</td>
					</tr>
				</tbody>
			</table>
		</div>
	`
};
