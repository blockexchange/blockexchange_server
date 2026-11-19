<script setup lang="ts">
import { computed, ref } from 'vue';
import PageHeader from "../PageHeader.vue";
import PagedTable from "../PagedTable.vue";
import { START, MEDIAFILES } from "../../util/breadcrumbs";
import { use_query_filter } from "../../util/query_filter";
import { has_permission } from "../../service/login";
import { get_mediafiles, count_mediafiles, get_mediafile_raw_url, delete_mediafile } from "../../api/media";
import type { MediafileInfo } from "../../api/types";
import format_size from "../../util/format_size";

const breadcrumb = [START, MEDIAFILES];
const allowed = computed(() => has_permission("MEDIA"));

const { filter, get_query } = use_query_filter();
const table = ref<{ update: () => Promise<void> } | null>(null);
const total = ref<number | null>(null);
const error = ref("");

const fetch_entries = (limit: number, offset: number) => get_mediafiles({ q: get_query(), limit, offset });
const count_entries = async () => (total.value = await count_mediafiles(get_query()));

const is_image = (name: string) => /\.(png|jpe?g|gif|bmp|tga)$/i.test(name);

async function remove(mf: MediafileInfo) {
	if (!confirm(`Delete the media file "${mf.name}"?`)) {
		return;
	}
	error.value = "";
	try {
		await delete_mediafile(mf.name);
		await table.value?.update();
	} catch (e) {
		error.value = `Deleting "${mf.name}" failed: ${JSON.stringify(e)}`;
	}
}
</script>

<template>
	<page-header title="Media files" subtitle="Textures and other media uploaded by the in-game mod, used to render schematic previews." :breadcrumb="breadcrumb"/>
	<div class="alert alert-warning" v-if="!allowed">
		<i class="fa fa-lock"></i> You need the <b>MEDIA</b> permission to manage media files.
	</div>
	<template v-else>
		<div class="alert alert-danger" v-if="error">{{ error }}</div>
		<div class="d-flex flex-wrap gap-2 align-items-center mb-3">
			<input class="form-control" style="max-width: 24rem;" placeholder="Filter by file or mod name" v-model="filter">
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
						<th style="width: 48px;"></th>
						<th>Name</th>
						<th class="d-none d-sm-table-cell">Mod</th>
						<th class="d-none d-md-table-cell text-end">Size</th>
						<th></th>
					</tr>
				</template>
				<template #body="{ list }">
					<tr v-if="total == 0">
						<td colspan="5" class="text-center text-muted-bx py-4">No media files found</td>
					</tr>
					<tr v-for="mf in list" :key="mf.name">
						<td>
							<img v-if="is_image(mf.name)" :src="get_mediafile_raw_url(mf.name)" width="32" height="32" alt="" loading="lazy" style="image-rendering: pixelated;"/>
							<i v-else class="fa fa-file fa-fw text-muted-bx"></i>
						</td>
						<td>
							<a class="font-monospace text-decoration-none" :href="get_mediafile_raw_url(mf.name)" target="_blank">{{ mf.name }}</a>
						</td>
						<td class="d-none d-sm-table-cell small text-muted-bx">{{ mf.mod_name }}</td>
						<td class="d-none d-md-table-cell small text-muted-bx text-end">{{ format_size(mf.size) }}</td>
						<td class="text-end">
							<button type="button" class="btn btn-sm btn-outline-danger" @click="remove(mf)" :title="'Delete ' + mf.name">
								<i class="fa fa-trash"></i>
								<span class="visually-hidden">Delete</span>
							</button>
						</td>
					</tr>
				</template>
			</paged-table>
		</div>
	</template>
</template>
