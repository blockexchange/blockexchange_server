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

onMounted(() => {
	get_collections_by_username(props.username).then(c => collections.value = c);
});
</script>

<template>
	<div v-if="!edit_mode">
		<router-link :to="'/collections/' + username + '/' + collection_name" class="btn btn-success w-100" v-if="collection_name">
			<i class="fa fa-object-group"></i>
			Part of the <span class="badge bg-primary">{{ collection_name }}</span> collection
		</router-link>
	</div>
	<div v-else>
		<select class="form-control" v-model="collection_uid">
			<option :value="null"></option>
			<option v-for="c in collections" :key="c.uid" :value="c.uid">{{ c.name }}</option>
		</select>
	</div>
</template>
