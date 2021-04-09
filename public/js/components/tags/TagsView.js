import { get_all, create, remove } from '../../api/tag.js';

import TagStore from '../../store/tag.js';
import Tag from './Tag.js';

const ListRow = {
	components: {
		"tag-label": Tag
	},
	props: ["tag"],
	methods: {
		remove: function(){
			remove(this.tag.id)
			.then(() => this.$emit("removed", this.tag.id));
		}
	},
	template: /*html*/`
	<tr>
		<td>
			<tag-label :tag="tag"/>
		</td>
		<td>{{ tag.description }}</td>
		<td>
			<button class="btn btn-danger" v-on:click="remove">
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
			.then(tag => {
				this.$emit("added", tag);
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
		"list-row": ListRow,
		"add-row": AddRow
	},
	data: function () {
		return TagStore;
	},
	methods: {
		update: function(){
			get_all()
			.then(tags => this.tags = tags);
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
				<list-row v-for="tag in tags"
					:key="tag.id"
					:tag="tag"
					v-on:removed="update"/>
				<add-row v-on:added="update"/>
			</tbody>
		</table>
	</div>
	`
};
