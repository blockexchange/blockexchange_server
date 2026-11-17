<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

import format_size from "../../util/format_size";
import format_time from "../../util/format_time";

import ModalPrompt from "../ModalPrompt.vue";
import SchemaStar from "./SchemaStar.vue";
import SchemaDownload from "./SchemaDownload.vue";
import SchemaCollection from "./SchemaCollection.vue";
import ModBadge from "../ModBadge.vue";
import TagBadge from "../TagBadge.vue";

import {
	schema_update,
	schema_set_tags,
	schema_update_screenshot,
	schema_delete,
	schema_update_mods,
	schema_update_info
} from "../../api/schema";
import type { SchemaSearchResponse, SchemaUpdateError } from "../../api/types";
import { get_tags as get_all_tags } from "../../service/tags";
import { has_permission } from "../../service/login";

const props = withDefaults(defineProps<{
	search_result: SchemaSearchResponse;
	allow_edit?: boolean;
}>(), {
	allow_edit: false
});

const emit = defineEmits<{
	save: [];
}>();

const router = useRouter();

const cache_counter = ref(Date.now());
const edit_mode = ref(false);
const error_response = ref<SchemaUpdateError | null>(null);
const screenshot_busy = ref(false);
const mods_busy = ref(false);
const delete_prompt = ref(false);

// the schema object is edited in place (reactive object owned by the parent's loading-block)
const schema = computed(() => props.search_result.schema);
const username = computed(() => props.search_result.username);
const mods = ref<string[]>(props.search_result.mods);
const tags = ref<string[]>(props.search_result.tags);

watch(() => props.search_result, sr => {
	mods.value = sr.mods;
	tags.value = sr.tags;
});

const available_tags = computed(() => get_all_tags().filter(t => !t.restricted || has_permission("ADMIN")));

function save() {
	error_response.value = null;
	schema_update(schema.value)
	.then(() => {
		// set tags
		schema_set_tags(schema.value.uid, tags.value);

		edit_mode.value = false;
		router.push(`/schema/${username.value}/${schema.value.name}`);
		emit("save");
	})
	.catch((e: SchemaUpdateError) => {
		error_response.value = e;
	});
}

function update_mods() {
	mods_busy.value = true;
	schema_update_mods(schema.value.uid)
	.then(m => {
		mods.value = m;
		mods_busy.value = false;
	});
}

function update_info() {
	schema_update_info(schema.value.uid)
	.then(s => {
		schema.value.total_size = s.total_size;
		schema.value.total_parts = s.total_parts;
	});
}

function update_screenshot() {
	screenshot_busy.value = true;
	schema_update_screenshot(schema.value.uid)
	.then(() => {
		screenshot_busy.value = false;
		cache_counter.value = Date.now();
	});
}

function delete_schema() {
	schema_delete(schema.value.uid)
	.then(() => router.push(`/schema/${username.value}`));
}

const markdown = (txt: string) => DOMPurify.sanitize(marked.parse(txt, { async: false }) as string);
</script>

<template>
	<div>
		<modal-prompt v-if="delete_prompt" title="Confirm deletion" @close="delete_prompt = false">
			<template #body>
				<p>Confirm deletion of schematic <b>{{ schema.name }}</b></p>
			</template>
			<template #footer>
				<a class="btn btn-secondary" @click="delete_prompt = false">
					<i class="fa fa-times"></i> Abort
				</a>
				<a class="btn btn-danger" @click="delete_schema">
					<i class="fa fa-trash"></i> Delete
				</a>
			</template>
		</modal-prompt>
		<div class="row">
			<div class="col-md-6">
				<h3 v-if="!edit_mode">
					{{ schema.name }}
					<small class="text-muted">by {{ username }}</small>
					&nbsp;
					<schema-star :schema_uid="schema.uid" v-model:stars="schema.stars"/>
				</h3>
				<div v-else class="input-group has-validation">
					<input type="text" class="form-control" v-model="schema.name" :class="{'is-invalid': error_response}">
					<div class="invalid-feedback" v-if="error_response && error_response.name_taken">
						Schematic name already taken
					</div>
					<div class="invalid-feedback" v-if="error_response && error_response.name_invalid">
						Schematic name invalid, allowed chars: a to z, A to Z, 0 to 9 and -, _
					</div>
				</div>
			</div>
			<div class="col-md-6">
				<div class="btn-group float-end" v-if="allow_edit && !edit_mode">
					<a class="btn btn-sm btn-secondary" @click="edit_mode = true">
						<i class="fa fa-edit"></i> Edit
					</a>
					<a class="btn btn-sm btn-secondary" :class="{'disabled': screenshot_busy}" @click="update_screenshot">
						<i v-if="screenshot_busy" class="fa fa-spinner fa-spin"></i>
						<i v-else class="fa fa-image"></i>
						Update screenshot
					</a>
					<a class="btn btn-sm btn-secondary" :class="{'disabled': mods_busy}" @click="update_mods">
						<i v-if="mods_busy" class="fa fa-spinner fa-spin"></i>
						<i v-else class="fa fa-box-archive"></i>
						Update modnames
					</a>
					<a class="btn btn-sm btn-secondary" @click="update_info">
						<i class="fa fa-chart-simple"></i>
						Update stats
					</a>
					<a class="btn btn-sm btn-danger" @click="delete_prompt = true">
						<i class="fa fa-trash"></i> Delete
					</a>
				</div>
				<div class="btn-group float-end" v-if="edit_mode">
					<a class="btn btn-sm btn-success" @click="save">
						<i class="fa fa-save"></i> Save
					</a>
				</div>
			</div>
		</div>
		<hr>
		<div class="row">
			<div class="col-md-4">
				<div class="card">
					<div class="card-header">
						<i class="fa fa-magnifying-glass"></i>
						Details
					</div>
					<div class="card-body">
						<ul>
							<li><b>Created:</b> {{ format_time(schema.created) }}</li>
							<li><b>Modified:</b> {{ format_time(schema.mtime) }}</li>
							<li><b>Size:</b> {{ format_size(schema.total_size) }}</li>
							<li><b>Dimensions:</b> {{ schema.size_x }} / {{ schema.size_y }} / {{ schema.size_z }} nodes</li>
							<li><b>Parts:</b> {{ schema.total_parts }}</li>
							<li><b>Downloads:</b> {{ schema.downloads }}</li>
							<li>
								<b>License:</b>
								<span v-if="!edit_mode">
									<img v-if="schema.license == 'CC0'" src="/pics/license_cc0.png">
									<img v-else-if="schema.license == 'CC-BY-SA'" src="/pics/license_cc-by-sa.png">
									<span v-else class="badge bg-secondary">{{ schema.license }}</span>
								</span>
								<select v-else class="form-control" v-model="schema.license">
									<option value="CC0">CC0</option>
									<option value="CC-BY-SA">CC-BY-SA</option>
									<option value="MIT">MIT</option>
									<option value="Proprietary">Proprietary</option>
								</select>
							</li>
						</ul>
					</div>
				</div>
				<br>
				<div class="card">
					<div class="card-header">
						<i class="fa fa-file-lines"></i>
						Description
					</div>
					<div class="card-body">
						<p v-if="!edit_mode">{{ schema.short_description }}</p>
						<input v-else type="text" class="form-control" v-model="schema.short_description" placeholder="Short description"/>
						<hr>
						<div v-if="!edit_mode" v-html="markdown(schema.description)"></div>
						<textarea v-else v-model="schema.description" rows="10" class="form-control"></textarea>
					</div>
				</div>
				<br>
				<div class="card">
					<div class="card-header">
						<i class="fa fa-box-archive"></i>
						Used mods
					</div>
					<div class="card-body">
						<a v-if="!edit_mode && schema.cdb_collection"
							class="btn btn-xs btn-outline-success w-100"
							target="new"
							:href="'https://content.minetest.net/collections/' + schema.cdb_collection + '/'">
							<i class="fa fa-box-archive"></i>
							Open CDB Mod collection
							<i class="fa fa-up-right-from-square"></i>
						</a>
						<hr v-if="schema.cdb_collection">
						<input v-if="edit_mode" type="text" class="form-control" v-model="schema.cdb_collection" placeholder="CDB Collection in the 'username/collectionname' form"/>
						<mod-badge v-for="mod in mods" :key="mod" :name="mod"/>
					</div>
				</div>
				<br>
				<div class="card">
					<div class="card-header">
						<i class="fa fa-tags"></i>
						Tags
					</div>
					<div class="card-body">
						<template v-if="!edit_mode">
							<tag-badge v-for="tag in tags" :key="tag" :name="tag"/>
						</template>
						<div v-else>
							<ul>
								<li v-for="tag in available_tags" :key="tag.uid">
									<input type="checkbox" class="form-check-input" :value="tag.name" v-model="tags"/>
									<span class="badge bg-success">
										<i class="fa fa-tag"></i>
										{{ tag.name }}
									</span>
								</li>
							</ul>
						</div>
					</div>
				</div>
				<br>
				<div class="card" v-if="edit_mode || search_result.collection_name">
					<div class="card-header">
						<i class="fa fa-object-group"></i>
						Collection
					</div>
					<div class="card-body">
						<schema-collection
							v-model="schema.collection_uid"
							:username="username"
							:edit_mode="edit_mode"
							:collection_name="search_result.collection_name"/>
					</div>
				</div>
			</div>
			<div class="col-md-8">
				<div class="card">
					<div class="card-header">
						<i class="fa fa-image"></i>
						Preview
					</div>
					<div class="card-body">
						<div class="text-center" style="min-height: 600px;">
							<img :src="'/api/schema/' + schema.uid + '/screenshot?cachebust=' + cache_counter" class="img-fluid">
						</div>
					</div>
				</div>
				<br>
				<div class="card">
					<div class="card-header">
						<i class="fa fa-download"></i>
						Download
					</div>
					<div class="card-body">
						<schema-download :schema="schema" :username="username"/>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>
