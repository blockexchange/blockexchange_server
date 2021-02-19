import { get_all, create, remove } from '../../api/access_token.js';
import store from '../../store/login.js';

const CreateForm = {
	data: function(){
		return {
			name: "",
			expires: 2,
			modifier: "86400000"
		};
	},
	methods: {
		add: function(){
			console.log(this.expires, this.modifier);
			create(this.name, Date.now() + (+this.expires * +this.modifier))
			.then(() => this.$emit('updated'));
		}
	},
	template: /*html*/`
	<div>
		<input type="text"
			class="form-control"
			v-model="name"
			placeholder="Token name (usually the name of the server you are using it on)"/>
		<input type="text"
			class="form-control"
			v-model="expires"
			placeholder="Expire-time"/>
		<select class="form-control" v-model="modifier">
			<option :value="3600*1000">Hours</option>
			<option :value="3600*1000*24">Days</option>
			<option :value="3600*1000*24*30">Months</option>
			<option :value="3600*1000*24*365">Years</option>
		</select>
		<button class="btn btn-primary" v-on:click="add" v-bind:disabled="!name">
			<i class="fa fa-plus"/> Add
		</button>
	</div>
	`
};

const List = {
	props: ["list"],
	data: function() {
		return {
			store: store
		};
	},
	methods: {
		remove: function(id){
			remove(id)
			.then(() => this.$emit("updated"));
		}
	},
	template: /*html*/`
	<table class="table table-condensed table-striped">
		<thead>
			<tr>
				<th>Name</th>
				<th>Created</th>
				<th>Expires</th>
				<th>Usage</th>
				<th>Use-count</th>
				<th>Action</th>
			</tr>
		</thead>
		<tbody>
			<tr v-for="token in list">
				<td>{{ token.name }}</td>
				<td>{{ new Date(+token.created).toLocaleString() }}</td>
				<td>{{ new Date(+token.expires).toLocaleString() }}</td>
				<td><pre>/bx_login {{store.claims.username}} {{ token.token }}</pre></td>
				<td>{{ token.usecount }}</td>
				<td>
					<button class="btn btn-sm btn-danger" v-on:click="remove(token.id)">
						<i class="fa fa-times"/> Remove
					</button>
				</td>
			</tr>
		</tbody>
	</table>
	`
};

export default {
	components: {
		"create-form": CreateForm,
		"token-list": List
	},
	data: function(){
		return {
			list: []
		};
	},
	methods: {
		update: function(){
			get_all().then(l => this.list = l);
		}
	},
	mounted: function(){
		this.update();
	},
	template: /*html*/`
	<div>
		<div class="row">
			<div class="col-md-12">
				<div class="card">
					<div class="card-header">
						Access token
					</div>
					<div class="card-body">
						<token-list :list="list" v-on:updated="update"/>
						<create-form v-on:updated="update"/>
					</div>
				</div>
			</div>
		</div>
	`
};
