<script setup lang="ts">
import BreadCrumb from "../Breadcrumb.vue";
import PagedTable from "../PagedTable.vue";
import { START, USERS } from "../../util/breadcrumbs";
import format_time from "../../util/format_time";
import { count_users, search_users } from "../../api/user";

const breadcrumb = [START, USERS];

const fetch_entries = (limit: number, offset: number) => search_users({ limit, offset });
const count_entries = () => count_users();
</script>

<template>
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
			<tr v-for="entry in list" :key="entry.uid">
				<td>
					<img v-if="entry.avatar_url" :src="entry.avatar_url" height="32" width="32" style="margin-right: 10px;"/>
					<router-link
						class="link-light link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover"
						:to="'/user/' + entry.name">
						<i class="fa fa-user"></i>
						{{ entry.name }}
					</router-link>
				</td>
				<td class="d-none d-sm-table-cell">
					{{ format_time(entry.created) }}
				</td>
				<td class="d-none d-sm-table-cell">
					<span class="badge bg-secondary">{{ entry.type }}</span>
				</td>
				<td class="d-none d-sm-table-cell">
					<span class="badge bg-secondary">{{ entry.role }}</span>
				</td>
			</tr>
		</template>
	</paged-table>
</template>
