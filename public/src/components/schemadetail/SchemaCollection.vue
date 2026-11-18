<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { get_collections_by_username } from "../../api/collection";
import type { Collection } from "../../api/types";

const props = defineProps<{
	username: string;
	edit_mode: boolean;
	collection_name: string | null;
}>();

// the collection_uid of the schema
const collection_uid = defineModel<string | null>({ required: true });

const collections = ref<Collection[]>([]);

onMounted(async () => {
	collections.value = await get_collections_by_username(props.username);
});
</script>

<template>
	<router-link v-if="!edit_mode && collection_name" :to="'/collections/' + username + '/' + collection_name" class="d-flex justify-content-between align-items-center text-decoration-none">
		<span><i class="fa fa-object-group me-2"></i>{{ collection_name }}</span>
		<span class="small">View all <i class="fa fa-chevron-right"></i></span>
	</router-link>
	<template v-if="edit_mode">
		<select class="form-select" v-model="collection_uid" aria-label="Collection">
			<option :value="null">Not in a collection</option>
			<option v-for="c in collections" :key="c.uid" :value="c.uid">{{ c.name }}</option>
		</select>
		<div class="form-text" v-if="collections.length == 0">
			You don't have any collections yet.
		</div>
	</template>
</template>
