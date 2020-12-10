export default {
	props: ["list"],
	template: /*html*/`
		<table class="table table-striped table-condensed">
			<thead>
				<tr>
					<th>Name</th>
					<th>User</th>
					<th>description</th>
					<th>Parts</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="entry in list">
					<td>
						<router-link :to="{ name: 'schemapage', params: { schemaname: entry.name, username: entry.user.name }}">
							{{ entry.name }}
						</router-link>
					</td>
					<td>{{ entry.user.name }}</td>
					<td>{{ entry.description }}</td>
					<td>{{ entry.total_parts }}</td>
				</tr>
			</tbody>
		</table>
	`
};
