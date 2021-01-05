import { get_by_userid } from '../../api/collection.js';
import store from '../../store/login.js';

const List = {
	created: function(){
		get_by_userid(store.claims.user_id)
		.then(c => this.collections = c);
	},
	data: function(){
		return {
			collections: []
		};
	},
	template: /*html*/`
		<div>
			<table class="table table-condensed table-striped">
				<thead>
					<tr>
						<th></th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td></td>
					</tr>
				</tbody>
			</table>
		</div>
	`
};

const Add = {
	template: /*html*/`
		<div>
			<input type="text"
				class="form-control"
				placeholder="Collection-name">
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
