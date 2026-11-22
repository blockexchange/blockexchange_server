<script setup lang="ts">
import { computed, ref } from 'vue';
import PageHeader from "../PageHeader.vue";
import EmptyState from "../EmptyState.vue";
import { START, TAGS } from "../../util/breadcrumbs";
import { get_tags, fetch_tags } from "../../service/tags";
import { has_permission } from "../../service/login";
import { create_tag, update_tag } from "../../api/tags";
import type { Tag, TagRequest } from "../../api/types";

const breadcrumb = [START, TAGS];
const tags = computed(get_tags);
const is_admin = computed(() => has_permission("ADMIN"));

const error = ref("");
const busy = ref(false);
const form_card = ref<HTMLElement | null>(null);

// uid of the tag being edited, null while adding a new one
const editing_uid = ref<string | null>(null);
const empty_form = (): TagRequest => ({ name: "", description: "", restricted: false });
const form = ref<TagRequest>(empty_form());

const name_taken = computed(() => tags.value.some(t => t.name == form.value.name && t.uid != editing_uid.value));
const form_valid = computed(() => form.value.name != "" && form.value.name.length <= 128 && !name_taken.value);

const error_message = (e: unknown) => (e as { message?: string })?.message || JSON.stringify(e);

function edit(tag: Tag) {
	editing_uid.value = tag.uid;
	form.value = { name: tag.name, description: tag.description, restricted: tag.restricted };
	error.value = "";
	form_card.value?.scrollIntoView({ behavior: "smooth", block: "nearest" });
}

function reset() {
	editing_uid.value = null;
	form.value = empty_form();
}

async function save() {
	busy.value = true;
	error.value = "";
	try {
		if (editing_uid.value) {
			await update_tag(editing_uid.value, form.value);
		} else {
			await create_tag(form.value);
		}
		reset();
		// refresh the shared tag list, it is used on the schema pages too
		await fetch_tags();
	} catch (e) {
		error.value = `Saving failed: ${error_message(e)}`;
	} finally {
		busy.value = false;
	}
}
</script>

<template>
	<page-header title="Tags" subtitle="Tags help people find schematics. Restricted tags can only be set by admins." :breadcrumb="breadcrumb"/>

	<template v-if="is_admin">
		<div class="alert alert-danger" v-if="error">{{ error }}</div>
		<div class="card mb-4" ref="form_card">
			<div class="card-body">
				<h6 class="card-title">
					<template v-if="editing_uid">
						<i class="fa fa-pen-to-square"></i> Edit tag
					</template>
					<template v-else>
						<i class="fa fa-plus"></i> Add tag
					</template>
				</h6>
				<form class="row g-2 align-items-end" @submit.prevent="save">
					<div class="col-sm-4 col-lg-3">
						<label class="form-label small" for="tag-name">Name</label>
						<input id="tag-name" class="form-control" maxlength="128" placeholder="e.g. medieval"
							:class="{'is-invalid': name_taken}" v-model.trim="form.name">
						<div class="invalid-feedback">A tag with this name already exists</div>
					</div>
					<div class="col-sm-8 col-lg-6">
						<label class="form-label small" for="tag-description">Description</label>
						<input id="tag-description" class="form-control" v-model.trim="form.description">
					</div>
					<div class="col-lg-3">
						<div class="form-check form-switch mb-2">
							<input id="tag-restricted" class="form-check-input" type="checkbox" role="switch" v-model="form.restricted">
							<label class="form-check-label" for="tag-restricted">
								<i class="fa fa-lock"></i> Restricted
							</label>
						</div>
					</div>
					<div class="col-12 d-flex gap-2 justify-content-end">
						<button class="btn btn-outline-secondary" type="button" v-if="editing_uid" @click="reset">
							Cancel
						</button>
						<button class="btn btn-success" type="submit" :disabled="!form_valid || busy">
							<span class="spinner-border spinner-border-sm" v-if="busy"></span>
							<i class="fa fa-floppy-disk" v-else></i>
							{{ editing_uid ? "Save" : "Add" }}
						</button>
					</div>
				</form>
			</div>
		</div>
	</template>

	<div class="list-group" v-if="tags.length">
		<div
			v-for="tag in tags"
			:key="tag.uid"
			class="list-group-item list-group-item-action position-relative"
			:class="{'active-tag': editing_uid == tag.uid}">
			<div class="d-flex justify-content-between align-items-center gap-2">
				<router-link :to="{path: '/search', query: {tag_name: tag.name}}" class="stretched-link text-reset text-decoration-none">
					<span class="badge text-bg-success me-2">{{ tag.name }}</span>{{ tag.description }}
				</router-link>
				<span class="d-flex align-items-center gap-2 text-nowrap">
					<span class="badge text-bg-secondary" v-if="tag.restricted"><i class="fa fa-lock"></i> restricted</span>
					<button type="button" class="btn btn-sm btn-outline-primary position-relative z-2"
						v-if="is_admin" @click="edit(tag)" :title="'Edit ' + tag.name">
						<i class="fa fa-pen"></i>
						<span class="visually-hidden">Edit</span>
					</button>
				</span>
			</div>
		</div>
	</div>
	<empty-state icon="tags" title="No tags defined" v-else/>
</template>

<style scoped>
.active-tag {
	background-color: var(--bs-tertiary-bg);
}
</style>
