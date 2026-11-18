<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { create_collection as api_create_collection, get_collections_by_username } from "../../api/collection";
import type { Collection } from "../../api/types";
import { get_user_uid, get_username } from "../../service/login";

import PageHeader from "../PageHeader.vue";
import EmptyState from "../EmptyState.vue";
import LoadingSpinner from "../LoadingSpinner.vue";
import { START, USER_COLLECTIONS } from "../../util/breadcrumbs";

const props = defineProps<{
	username: string;
	readonly?: boolean;
}>();

const breadcrumb = computed(() => [START, USER_COLLECTIONS(props.username)]);
const collections = ref<Collection[]>([]);
const busy = ref(false);
const can_create = computed(() => !props.readonly && get_username() == props.username);
const new_collection_name = ref("");

const valid_name = computed(() => /^[a-zA-Z0-9_.-]*$/.test(new_collection_name.value));
const duplicate_name = computed(() => collections.value.some(c => c.name == new_collection_name.value));

async function update() {
	busy.value = true;
	try {
		const c = await get_collections_by_username(props.username);
		c.sort((a, b) => a.name.localeCompare(b.name));
		collections.value = c;
	} finally {
		busy.value = false;
	}
}

async function create_collection() {
	await api_create_collection({
		name: new_collection_name.value,
		user_uid: get_user_uid() ?? ""
	});
	new_collection_name.value = "";
	await update();
}

onMounted(update);
</script>

<template>
	<page-header :title="'Collections by ' + username" :breadcrumb="breadcrumb"
		subtitle="Collections group related schematics, such as all parts of a city or a themed build set."/>
	<loading-spinner v-if="busy"/>
	<div class="row g-4" v-else>
		<div class="col-lg-8">
			<div class="list-group" v-if="collections.length">
				<router-link
					v-for="collection in collections"
					:key="collection.uid"
					class="list-group-item list-group-item-action d-flex justify-content-between align-items-center"
					:to="'/collections/' + username + '/' + collection.name">
					<span><i class="fa fa-object-group me-2 text-muted-bx"></i>{{ collection.name }}</span>
					<i class="fa fa-chevron-right text-muted-bx"></i>
				</router-link>
			</div>
			<empty-state icon="object-group" title="No collections yet" v-else>
				<template v-if="can_create">Create your first collection with the form.</template>
			</empty-state>
		</div>
		<div class="col-lg-4" v-if="can_create">
			<form class="card" @submit.prevent="create_collection">
				<div class="card-header">New collection</div>
				<div class="card-body">
					<label class="form-label small" for="collection-name">Name</label>
					<input id="collection-name" type="text"
						class="form-control"
						:class="{'is-invalid': !valid_name || duplicate_name}"
						placeholder="e.g. medieval_town"
						v-model="new_collection_name"/>
					<div class="invalid-feedback" v-if="new_collection_name && !valid_name">
						Only letters (a-z, A-Z), numbers and - _ . are allowed.
					</div>
					<div class="invalid-feedback" v-if="new_collection_name && duplicate_name">
						You already have a collection with this name.
					</div>
					<button class="btn btn-success w-100 mt-3" type="submit"
						:disabled="!valid_name || duplicate_name || !new_collection_name">
						<i class="fa fa-plus"></i> Create collection
					</button>
				</div>
			</form>
		</div>
	</div>
</template>
