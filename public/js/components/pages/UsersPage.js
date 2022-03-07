import { get_all } from '../../api/user.js';

export default {
	created: function(){
		get_all().then(userdata => this.userdata = userdata);
	},
	data: function(){
		return {
			userdata: null
		};
	},
	template: /*html*/`
		<div v-if="userdata">
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
