<script setup lang="ts">
import { computed, ref } from 'vue';
import PageHeader from "../PageHeader.vue";
import PagedTable from "../PagedTable.vue";
import { START, NODEDEFINITIONS } from "../../util/breadcrumbs";
import { use_query_filter } from "../../util/query_filter";
import { has_permission } from "../../service/login";
import { get_nodedefinitions, count_nodedefinitions, delete_nodedefinition } from "../../api/media";
import type { Nodedefinition } from "../../api/types";

const breadcrumb = [START, NODEDEFINITIONS];
const allowed = computed(() => has_permission("MEDIA"));

const { filter, get_query } = use_query_filter();
const table = ref<{ update: () => Promise<void> } | null>(null);
const total = ref<number | null>(null);
const expanded = ref<string | null>(null);
const error = ref("");

const fetch_entries = (limit: number, offset: number) => get_nodedefinitions({ q: get_query(), limit, offset });
const count_entries = async () => (total.value = await count_nodedefinitions(get_query()));

async function remove(nd: Nodedefinition) {
	if (!confirm(`Delete the node definition "${nd.name}"?`)) {
		return;
	}
	error.value = "";
	try {
		await delete_nodedefinition(nd.name);
		await table.value?.update();
	} catch (e) {
		error.value = `Deleting "${nd.name}" failed: ${JSON.stringify(e)}`;
	}
}
</script>

<template>
	<page-header title="Node definitions" subtitle="Node definitions uploaded by the in-game mod, used to render schematic previews." :breadcrumb="breadcrumb"/>
	<div class="alert alert-warning" v-if="!allowed">
		<i class="fa fa-lock"></i> You need the <b>MEDIA</b> permission to manage node definitions.
	</div>
	<template v-else>
		<div class="alert alert-danger" v-if="error">{{ error }}</div>
		<div class="d-flex flex-wrap gap-2 align-items-center mb-3">
			<input class="form-control" style="max-width: 24rem;" placeholder="Filter by node or mod name" v-model="filter">
			<span class="text-muted-bx small" v-if="total != null">{{ total }} entries</span>
		</div>
		<div class="card">
			<paged-table
				ref="table"
				class="table table-hover align-middle mb-0"
				:per_page="100"
				:fetch_entries="fetch_entries"
				:count_entries="count_entries">
				<template #header>
					<tr>
						<th>Name</th>
						<th class="d-none d-sm-table-cell">Mod</th>
						<th></th>
					</tr>
				</template>
				<template #body="{ list }">
					<tr v-if="total == 0">
						<td colspan="3" class="text-center text-muted-bx py-4">No node definitions found</td>
					</tr>
					<template v-for="nd in list" :key="nd.name">
						<tr>
							<td>
								<button type="button" class="btn btn-link p-0 text-decoration-none font-monospace" @click="expanded = expanded == nd.name ? null : nd.name">
									<i class="fa fa-fw" :class="expanded == nd.name ? 'fa-caret-down' : 'fa-caret-right'"></i>
									{{ nd.name }}
								</button>
							</td>
							<td class="d-none d-sm-table-cell small text-muted-bx">{{ nd.mod_name }}</td>
							<td class="text-end">
								<button type="button" class="btn btn-sm btn-outline-danger" @click="remove(nd)" :title="'Delete ' + nd.name">
									<i class="fa fa-trash"></i>
									<span class="visually-hidden">Delete</span>
								</button>
							</td>
						</tr>
						<tr v-if="expanded == nd.name">
							<td colspan="3">
								<pre class="small mb-0" style="white-space: pre-wrap; word-break: break-all;">{{ nd.definition }}</pre>
							</td>
						</tr>
					</template>
				</template>
			</paged-table>
		</div>
	</template>
</template>
