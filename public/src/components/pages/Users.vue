<script setup lang="ts">
import PageHeader from "../PageHeader.vue";
import PagedTable from "../PagedTable.vue";
import { START, USERS } from "../../util/breadcrumbs";
import { count_users, search_users } from "../../api/user";

const breadcrumb = [START, USERS];

const fetch_entries = (limit: number, offset: number) => search_users({ limit, offset });
const count_entries = () => count_users();

const account_type = (t: string) => t == "LOCAL" ? "Password" : t.charAt(0) + t.slice(1).toLowerCase();
</script>

<template>
	<page-header title="Users" subtitle="Everyone who has an account here. Open a profile to see that user's schematics." :breadcrumb="breadcrumb"/>
	<div class="card">
		<paged-table
			class="table table-hover align-middle mb-0"
			:fetch_entries="fetch_entries"
			:count_entries="count_entries">
			<template #header>
				<tr>
					<th>Name</th>
					<th class="d-none d-sm-table-cell">Joined</th>
					<th class="d-none d-md-table-cell">Login</th>
				</tr>
			</template>
			<template #body="{ list }">
				<tr v-for="entry in list" :key="entry.uid">
					<td>
						<router-link class="link-light text-decoration-none d-inline-flex align-items-center gap-2" :to="'/user/' + entry.name">
							<img v-if="entry.avatar_url" :src="entry.avatar_url" height="32" width="32" class="rounded" alt="" loading="lazy"/>
							<span v-else class="rounded bg-secondary d-inline-flex align-items-center justify-content-center" style="width: 32px; height: 32px;">
								<i class="fa fa-user text-muted-bx"></i>
							</span>
							<span>{{ entry.name }}</span>
							<span class="badge bg-warning text-dark" v-if="entry.role == 'ADMIN'">Admin</span>
						</router-link>
					</td>
					<td class="d-none d-sm-table-cell small text-muted-bx">
						{{ new Date(entry.created).toLocaleDateString() }}
					</td>
					<td class="d-none d-md-table-cell small text-muted-bx">
						{{ account_type(entry.type) }}
					</td>
				</tr>
			</template>
		</paged-table>
	</div>
</template>
