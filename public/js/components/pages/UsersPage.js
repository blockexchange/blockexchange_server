import { get_all } from '../../api/user.js';

export default {
	created: function(){
		get_all().then(users => this.users = users);
	},
	data: function(){
		return {
			users: []
		};
	},
	template: /*html*/`
		<div>
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
					<tr v-for="user in users">
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
