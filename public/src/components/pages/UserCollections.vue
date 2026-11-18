<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';

import { create_collection as api_create_collection, get_collections_by_username } from "../../api/collection";
import type { Collection } from "../../api/types";
import { get_user_uid } from "../../service/login";

import BreadCrumb from "../Breadcrumb.vue";
import LoadingSpinner from "../LoadingSpinner.vue";
import { START, USER_COLLECTIONS } from "../../util/breadcrumbs";

const props = defineProps<{
	username: string;
	readonly?: boolean;
}>();

const breadcrumb = computed(() => [START, USER_COLLECTIONS(props.username)]);
const collections = ref<Collection[]>([]);
const busy = ref(false);
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
	<bread-crumb :items="breadcrumb"/>
	<loading-spinner v-if="busy"/>
	<div v-else>
		<table class="table table-dark table-striped table-condensed">
			<thead>
				<tr>
					<th>Name</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<tr :key="collection.uid" v-for="collection in collections">
					<td>
						<router-link
							class="link-info link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover"
							:to="'/collections/' + username + '/' + collection.name">
							<i class="fa fa-object-group"></i>
							{{ collection.name }}
						</router-link>
					</td>
					<td>
						<div class="btn-group">
							<button class="btn btn-danger">
								<i class="fa fa-trash"></i>
								Remove
							</button>
							<button class="btn btn-primary">
								<i class="fa fa-edit"></i>
								Edit
							</button>
						</div>
					</td>
				</tr>
				<tr v-if="!readonly">
					<td>
						<input type="text"
							class="form-control"
							:class="{'is-invalid': !valid_name || duplicate_name}"
							placeholder="New collection name"
							v-model="new_collection_name"/>
						<div class="invalid-feedback" v-if="new_collection_name && !valid_name">
							Name is invalid, allowed chars: a to z, A to Z, 0 to 9 and -._
						</div>
						<div class="invalid-feedback" v-if="new_collection_name && duplicate_name">
							Name is already taken
						</div>
					</td>
					<td>
						<button class="btn btn-secondary"
							:disabled="!valid_name || duplicate_name || !new_collection_name"
							@click="create_collection">
							<i class="fa fa-plus"></i> Create collection
						</button>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
</template>
