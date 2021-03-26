import { get_by_userid, create } from '../../api/collection.js';
import loginstore from '../../store/login.js';

const ListRows = {
	props: ["collections"],
	created: function(){
		console.log("ListRows::created", this.collections);
	},
	template: /*html*/`
	<tr v-for="collection in collections" :key="collection.id">
		<td>{{ collection.name }}</td>
		<td>{{ collection.description }}</td>
		<td>
			<button class="btn btn-danger">
				Delete
			</button>
		</td>
	</tr>
	`
};

const AddRow = {
	data: function () {
		return {
			new_name: "",
			new_description: ""
		};
	},
	methods: {
		create: function(){
			create({
				name: this.new_name,
				description: this.new_description
			})
			.then(collection => {
				this.$emit("added", collection);
			});
		}
	},
	template: /*html*/`
	<tr>
		<td>
			<input type="text" class="form-control" placeholder="Name" v-model="new_name"/>
		</td>
		<td>
			<input type="text" class="form-control" placeholder="Description" v-model="new_description"/>
		</td>
		<td>
			<button class="btn btn-success"
				v-bind:disabled="!new_name || !new_description"
				v-on:click="create">
				Add
			</button>
		</td>
	</tr>
	`
};


export default {
	components: {
		"list-rows": ListRows,
		"add-row": AddRow
	},
	created: function () {
		this.update();
	},
	data: function () {
		return {
			collections: []
		};
	},
	methods: {
		update: function(){
			get_by_userid(loginstore.claims.user_id)
			.then(c => {
				this.$set(this, "collections", c);
				console.log(this.collections);
			});
		}
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
				<list-rows :collections="collections"></list-rows>
				<add-row v-on:added="update"/>
			</tbody>
		</table>
	</div>
	`
};
