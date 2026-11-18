<script setup lang="ts">
import { onMounted, ref } from 'vue';

import SchemaList from "./SchemaList.vue";
import { schema_search } from "../api/schema";
import type { SchemaSearchResponse } from "../api/types";
import { get_tags } from "../service/tags";

const list = ref<SchemaSearchResponse[]>([]);

onMounted(async () => {
	const featured = get_tags().find(t => t.name == "featured");
	if (!featured) {
		return;
	}

	list.value = await schema_search({
		tag_uid: featured.uid,
		limit: 12,
		complete: true
	});
});
</script>

<template>
	<schema-list :list="list" :show_not_found="false"/>
</template>
