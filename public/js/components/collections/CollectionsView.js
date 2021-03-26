import { get_by_userid } from '../../api/collection.js';
import store from '../../store/login.js';

const List = {
	created: function () {
		get_by_userid(store.claims.user_id)
			.then(c => this.collections = c);
	},
	data: function () {
		return {
			collections: [],
			new_name: "",
			new_description: ""
		};
	},
	template: /*html*/`
		<div>
			<table class="table table-condensed table-striped">
				<thead>
					<tr>
						<th>Name</th>
						<th>Description</th>
						<th>Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="collection in collections">
						<td>{{ collection.name }}</td>
						<td>{{ collection.description }}</td>
						<td>
							<button class="btn btn-danger">
								Delete
							</button>
						</td>
					</tr>
					<tr>
						<td>
							<input type="text" class="form-control" placeholder="Name" v-model="new_name"/>
						</td>
						<td>
							<input type="text" class="form-control" placeholder="Description" v-model="new_description"/>
						</td>
						<td>
							<button class="btn btn-success" v-bind:disabled="!new_name || !new_description">
								Add
							</button>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
	`
};

const Add = {
	template: /*html*/`
		<div>
		</div>
	`
};


export default {
	components: {
		"collection-list": List,
		"collection-add": Add
	},
	template: /*html*/`
	<div>
		<collection-list/>
		<collection-add/>
	</div>
	`
};
