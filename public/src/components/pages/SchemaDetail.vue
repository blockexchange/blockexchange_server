<script setup lang="ts">
import { computed } from 'vue';

import BreadCrumb from "../Breadcrumb.vue";
import LoadingBlock from "../LoadingBlock.vue";
import EmptyState from "../EmptyState.vue";
import SchemaDetail from "../schemadetail/SchemaDetail.vue";
import { START, USERS, USER_DETAIL, SCHEMA_DETAIL } from "../../util/breadcrumbs";

import { has_permission, get_user_uid } from "../../service/login";
import { schema_search } from "../../api/schema";
import type { Schema } from "../../api/types";

const props = defineProps<{
	username: string;
	name: string;
}>();

const breadcrumb = computed(() => [START, USERS, USER_DETAIL(props.username), SCHEMA_DETAIL(props.username, props.name)]);

const fetch_search_result = async () => {
	const r = await schema_search({
		schema_name: props.name,
		user_name: props.username
	});
	return r[0] ?? null;
};

const fetch_data = () => ({
	search_result: fetch_search_result()
});

const allow_edit = (schema: Schema) => (get_user_uid() == schema.user_uid || has_permission("ADMIN"));
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<loading-block :fetch_data="fetch_data" v-slot="{ data, update_data }">
		<empty-state icon="circle-question" title="Schematic not found" v-if="!data.search_result">
			<p><b>{{ name }}</b> by {{ username }} doesn't exist or was renamed.</p>
			<router-link :to="'/user/' + username">See schematics by {{ username }}</router-link>
		</empty-state>
		<schema-detail
			v-else
			:search_result="data.search_result"
			:allow_edit="allow_edit(data.search_result.schema)"
			@save="update_data"/>
	</loading-block>
</template>
