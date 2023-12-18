import Breadcrumb, { START, USERS } from "../Breadcrumb.js";
import PagedTable from "../PagedTable.js";
import format_time from "../../util/format_time.js";
import { count_users, search_users } from "../../api/user.js";

export default {
	components: {
        "bread-crumb": Breadcrumb,
		"paged-table": PagedTable
	},
	data: function() {
		return {
			breadcrumb: [START, USERS],
			users: []
		};
	},
	methods: {
		format_time,
		fetch_entries: function(limit, offset) {
			return search_users({ limit, offset });
		},
		count_entries: function() {
			return count_users();
		}
	},
	template: /*html*/`
		<bread-crumb :items="breadcrumb"/>
		<paged-table
			class="table table-dark table-condensed table-striped"
			:fetch_entries="fetch_entries"
			:count_entries="count_entries">
			<template #header>
				<tr>
					<th>Name</th>
					<th class="d-none d-sm-table-cell">Created</th>
					<th class="d-none d-sm-table-cell">Type</th>
					<th class="d-none d-sm-table-cell">Role</th>
				</tr>
			</template>
			<template #body="{ list }">
				<tr v-for="entry in list">
					<td>
						<router-link :to="'/user/' + entry.name">
							{{entry.name}}
						</router-link>
					</td>
					<td class="d-none d-sm-table-cell">
						{{format_time(entry.created)}}
					</td>
					<td class="d-none d-sm-table-cell">
						<span class="badge bg-secondary">{{entry.type}}</span>
					</td>
					<td class="d-none d-sm-table-cell">
						<span class="badge bg-secondary">{{entry.role}}</span>
					</td>
				</tr>
			</template>
		</paged-table>
	`
};
