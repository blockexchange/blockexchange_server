<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRouter } from 'vue-router';
import { marked } from 'marked';
import DOMPurify from 'dompurify';

import format_size from "../../util/format_size";
const format_date = (ts: number) => new Date(ts).toLocaleDateString();

import ModalPrompt from "../ModalPrompt.vue";
import SchemaStar from "./SchemaStar.vue";
import SchemaDownload from "./SchemaDownload.vue";
import SchemaCollection from "./SchemaCollection.vue";
import SchemaModDependencies from "./SchemaModDependencies.vue";
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
const tools_open = ref(false);
const saving = ref(false);

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

async function save() {
	error_response.value = null;
	saving.value = true;
	try {
		await schema_update(schema.value);
		// set tags
		await schema_set_tags(schema.value.uid, tags.value);
	} catch (e) {
		error_response.value = e as SchemaUpdateError;
		return;
	} finally {
		saving.value = false;
	}

	edit_mode.value = false;
	await router.push(`/schema/${username.value}/${schema.value.name}`);
	emit("save");
}

function cancel_edit() {
	// discard the in-place edits by reloading the schematic
	edit_mode.value = false;
	error_response.value = null;
	emit("save");
}

function start_edit() {
	tools_open.value = false;
	edit_mode.value = true;
}

async function update_mods() {
	tools_open.value = false;
	mods_busy.value = true;
	try {
		mods.value = await schema_update_mods(schema.value.uid);
	} finally {
		mods_busy.value = false;
	}
}

async function update_info() {
	tools_open.value = false;
	const s = await schema_update_info(schema.value.uid);
	schema.value.total_size = s.total_size;
	schema.value.total_parts = s.total_parts;
}

async function update_screenshot() {
	tools_open.value = false;
	screenshot_busy.value = true;
	try {
		await schema_update_screenshot(schema.value.uid);
		cache_counter.value = Date.now();
	} finally {
		screenshot_busy.value = false;
	}
}

async function delete_schema() {
	await schema_delete(schema.value.uid);
	await router.push(`/schema/${username.value}`);
}

const markdown = (txt: string) => DOMPurify.sanitize(marked.parse(txt, { async: false }) as string);
</script>

<template>
	<div>
		<modal-prompt v-if="delete_prompt" title="Delete schematic?" @close="delete_prompt = false">
			<template #body>
				<p class="mb-0">
					<b>{{ schema.name }}</b> will be removed for everyone, including the in-game
					<code>/bx_load</code> command. This can't be undone.
				</p>
			</template>
			<template #footer>
				<button type="button" class="btn btn-secondary" @click="delete_prompt = false">Cancel</button>
				<button type="button" class="btn btn-danger" @click="delete_schema">
					<i class="fa fa-trash"></i> Delete schematic
				</button>
			</template>
		</modal-prompt>

		<!-- header -->
		<div class="d-flex flex-wrap justify-content-between align-items-start gap-3 mb-4">
			<div class="flex-grow-1" style="min-width: 0;">
				<template v-if="!edit_mode">
					<h1 class="h2 mb-1 text-break">{{ schema.name }}</h1>
					<div class="text-muted-bx">
						by <router-link :to="'/user/' + username">{{ username }}</router-link>
						<template v-if="schema.short_description"> · {{ schema.short_description }}</template>
					</div>
				</template>
				<div v-else class="row g-2">
					<div class="col-md-5">
						<label class="form-label small" for="schema-name">Name</label>
						<input id="schema-name" type="text" class="form-control" v-model="schema.name" :class="{'is-invalid': error_response}">
						<div class="invalid-feedback" v-if="error_response && error_response.name_taken">
							You already have a schematic with this name.
						</div>
						<div class="invalid-feedback" v-else-if="error_response && error_response.name_invalid">
							Only letters (a-z, A-Z), numbers, - and _ are allowed.
						</div>
						<div class="invalid-feedback" v-else>Saving failed.</div>
					</div>
					<div class="col-md-7">
						<label class="form-label small" for="schema-short">Short description</label>
						<input id="schema-short" type="text" class="form-control" v-model="schema.short_description" placeholder="One line, shown in search results"/>
					</div>
				</div>
			</div>
			<div class="d-flex flex-wrap gap-2 align-items-center">
				<template v-if="!edit_mode">
					<router-link :to="'/schema/' + username + '/' + schema.name + '/3d'" class="btn btn-outline-primary">
						<i class="fa fa-cube"></i> 3D View
					</router-link>
					<schema-star :schema_uid="schema.uid" v-model:stars="schema.stars"/>
					<template v-if="allow_edit">
						<button type="button" class="btn btn-outline-secondary" @click="start_edit">
							<i class="fa fa-pen"></i> Edit
						</button>
						<div class="dropdown">
							<button type="button" class="btn btn-outline-secondary dropdown-toggle" :aria-expanded="tools_open" @click="tools_open = !tools_open">
								<span class="spinner-border spinner-border-sm" v-if="screenshot_busy || mods_busy"></span>
								<i class="fa fa-screwdriver-wrench" v-else></i>
								<span class="visually-hidden">More actions</span>
							</button>
							<ul class="dropdown-menu dropdown-menu-end" :class="{show: tools_open}" style="right: 0;">
								<li>
									<button class="dropdown-item" type="button" :disabled="screenshot_busy" @click="update_screenshot">
										<i class="fa fa-image fa-fw"></i> Regenerate preview
									</button>
								</li>
								<li>
									<button class="dropdown-item" type="button" :disabled="mods_busy" @click="update_mods">
										<i class="fa fa-box-archive fa-fw"></i> Re-scan used mods
									</button>
								</li>
								<li>
									<button class="dropdown-item" type="button" @click="update_info">
										<i class="fa fa-chart-simple fa-fw"></i> Recalculate size
									</button>
								</li>
								<li><hr class="dropdown-divider"></li>
								<li>
									<button class="dropdown-item text-danger" type="button" @click="tools_open = false; delete_prompt = true">
										<i class="fa fa-trash fa-fw"></i> Delete
									</button>
								</li>
							</ul>
						</div>
					</template>
				</template>
				<template v-else>
					<button type="button" class="btn btn-secondary" @click="cancel_edit">Cancel</button>
					<button type="button" class="btn btn-success" :disabled="saving" @click="save">
						<span class="spinner-border spinner-border-sm" v-if="saving"></span>
						<i class="fa fa-check" v-else></i>
						Save changes
					</button>
				</template>
			</div>
		</div>

		<div class="row g-4">
			<!-- main column -->
			<div class="col-lg-8 d-flex flex-column gap-4">
				<div class="card overflow-hidden">
					<div class="text-center bx-preview">
						<img :src="'/api/schema/' + schema.uid + '/screenshot?cachebust=' + cache_counter" class="img-fluid" :alt="'Preview of ' + schema.name">
					</div>
				</div>

				<section class="card" v-if="edit_mode || schema.description">
					<div class="card-header">Description</div>
					<div class="card-body">
						<div v-if="!edit_mode" class="markdown-body" v-html="markdown(schema.description)"></div>
						<template v-else>
							<textarea v-model="schema.description" rows="10" class="form-control font-monospace"></textarea>
							<div class="form-text">Markdown is supported.</div>
						</template>
					</div>
				</section>

				<section class="card" v-if="edit_mode && allow_edit">
					<div class="card-header"><i class="fa fa-code-branch"></i> Mod dependencies</div>
					<div class="card-body">
						<schema-mod-dependencies :schema_uid="schema.uid"/>
					</div>
				</section>

				<section class="card" v-if="!edit_mode">
					<div class="card-header"><i class="fa fa-download"></i> Use this schematic</div>
					<div class="card-body">
						<schema-download :schema="schema" :username="username"/>
					</div>
				</section>
			</div>

			<!-- sidebar -->
			<aside class="col-lg-4 d-flex flex-column gap-4">
				<section class="card">
					<div class="card-header">Details</div>
					<div class="card-body">
						<dl class="bx-dl row mb-0">
							<dt class="col-5">Dimensions</dt>
							<dd class="col-7" title="Width × height × depth in nodes">{{ schema.size_x }} × {{ schema.size_y }} × {{ schema.size_z }}</dd>
							<dt class="col-5">Size</dt>
							<dd class="col-7">{{ format_size(schema.total_size) }} <span class="text-muted-bx small">({{ schema.total_parts }} parts)</span></dd>
							<dt class="col-5">Downloads</dt>
							<dd class="col-7">{{ schema.downloads }}</dd>
							<dt class="col-5">Created</dt>
							<dd class="col-7">{{ format_date(schema.created) }}</dd>
							<dt class="col-5">Updated</dt>
							<dd class="col-7">{{ format_date(schema.mtime) }}</dd>
							<dt class="col-5">License</dt>
							<dd class="col-7 mb-0">
								<template v-if="!edit_mode">
									<img v-if="schema.license == 'CC0'" src="/pics/license_cc0.png" alt="CC0" title="CC0: public domain">
									<img v-else-if="schema.license == 'CC-BY-SA'" src="/pics/license_cc-by-sa.png" alt="CC-BY-SA" title="CC-BY-SA: credit the author, share alike">
									<span v-else>{{ schema.license }}</span>
								</template>
								<select v-else class="form-select form-select-sm" v-model="schema.license" aria-label="License">
									<option value="CC0">CC0</option>
									<option value="CC-BY-SA">CC-BY-SA</option>
									<option value="MIT">MIT</option>
									<option value="Proprietary">Proprietary</option>
								</select>
							</dd>
						</dl>
					</div>
				</section>

				<section class="card">
					<div class="card-header">Required mods</div>
					<div class="card-body">
						<p class="small text-muted-bx" v-if="!edit_mode && mods.length">Install these in your world before placing the schematic.</p>
						<div class="d-flex flex-wrap gap-1">
							<mod-badge v-for="mod in mods" :key="mod" :name="mod"/>
							<span class="small text-muted-bx" v-if="!mods.length">No mods detected.</span>
						</div>
						<a v-if="!edit_mode && schema.cdb_collection"
							class="btn btn-sm btn-outline-success w-100 mt-3"
							target="_blank"
							rel="noopener"
							:href="'https://content.luanti.org/collections/' + schema.cdb_collection + '/'">
							Get all mods as a ContentDB collection
							<i class="fa fa-up-right-from-square"></i>
						</a>
						<div v-if="edit_mode" class="mt-3">
							<label class="form-label small" for="schema-cdb">ContentDB collection</label>
							<input id="schema-cdb" type="text" class="form-control form-control-sm" v-model="schema.cdb_collection" placeholder="username/collection_name"/>
						</div>
					</div>
				</section>

				<section class="card" v-if="edit_mode || tags.length">
					<div class="card-header">Tags</div>
					<div class="card-body">
						<div class="d-flex flex-wrap gap-1" v-if="!edit_mode">
							<tag-badge v-for="tag in tags" :key="tag" :name="tag"/>
						</div>
						<div v-else>
							<div class="form-check" v-for="tag in available_tags" :key="tag.uid">
								<input type="checkbox" class="form-check-input" :id="'tag-' + tag.uid" :value="tag.name" v-model="tags"/>
								<label class="form-check-label" :for="'tag-' + tag.uid">
									{{ tag.name }}
									<span class="small text-muted-bx" v-if="tag.description">· {{ tag.description }}</span>
								</label>
							</div>
						</div>
					</div>
				</section>

				<section class="card" v-if="edit_mode || search_result.collection_name">
					<div class="card-header">Collection</div>
					<div class="card-body">
						<schema-collection
							v-model="schema.collection_uid"
							:username="username"
							:edit_mode="edit_mode"
							:collection_name="search_result.collection_name"/>
					</div>
				</section>
			</aside>
		</div>
	</div>
</template>
