<script setup lang="ts">
import { onMounted, ref } from 'vue';

import SchemaList from "./SchemaList.vue";
import LoadingSpinner from "./LoadingSpinner.vue";
import EmptyState from "./EmptyState.vue";
import { schema_search } from "../api/schema";
import type { SchemaSearchResponse } from "../api/types";
import { get_tags } from "../service/tags";

const list = ref<SchemaSearchResponse[]>([]);
const busy = ref(true);

onMounted(async () => {
	try {
		const featured = get_tags().find(t => t.name == "featured");
		if (featured) {
			list.value = await schema_search({
				tag_uid: featured.uid,
				limit: 8,
				complete: true
			});
		}
	} finally {
		busy.value = false;
	}
});
</script>

<template>
	<loading-spinner v-if="busy"/>
	<schema-list :list="list" v-else>
		<template #empty>
			<empty-state icon="star" title="Nothing featured yet">
				<router-link to="/search">Browse all schematics</router-link> instead.
			</empty-state>
		</template>
	</schema-list>
</template>
