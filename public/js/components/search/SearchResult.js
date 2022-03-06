import store from '../../store/login.js';
import Tag from '../tags/Tag.js';
import prettysize from '../../util/prettysize.js';

export default {
	data: function(){
		return {
			store: store
		};
	},
	components: {
		"tag-label": Tag
	},
	props: ["list"],
	methods: {
		prettysize: prettysize
	},
	template: /*html*/`
		<table class="table table-striped table-condensed">
			<thead>
				<tr>
					<th>Name</th>
					<th>User</th>
					<th>description</th>
					<th>Parts</th>
					<th>Size</th>
					<th>Changed</th>
				</tr>
			</thead>
			<tbody>
				<tr v-for="entry in list" v-bind:class="{ 'table-danger': !entry.complete }">
					<td>
						<router-link :to="{ name: 'schemapage', params: { schema_name: entry.name, user_name: entry.user.name }}">
							{{ entry.name }}
						</router-link>
						<tag-label v-for="tag in entry.tags" :tag_id="tag.tag_id"/>
					</td>
					<td>
						{{ entry.user.name }}
						<span v-if="store.claims && store.claims.username == entry.user.name" class="badge bg-secondary">
							owner
						</span>
					</td>
					<td>{{ entry.description }}</td>
					<td>
						{{ entry.total_parts }}
						<span v-if="!entry.complete" class="badge bg-danger">
							Incomplete
						</span>
					</td>
					<td>{{ prettysize(entry.total_size) }}</td>
					<td>{{ new Date(+entry.created).toDateString() }}</td>
				</tr>
			</tbody>
		</table>
	`
};
