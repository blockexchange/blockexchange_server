<script setup lang="ts">
import { computed } from 'vue';

import BreadCrumb from "../Breadcrumb.vue";
import UserProfile from "../UserProfile.vue";
import PagedContent from "../PagedContent.vue";
import SchemaList from "../SchemaList.vue";
import EmptyState from "../EmptyState.vue";
import { START, USERS, USER_DETAIL } from "../../util/breadcrumbs";
import { schema_count, schema_search } from "../../api/schema";
import type { SchemaSearchRequest } from "../../api/types";
import { get_username, has_permission } from "../../service/login";

const props = defineProps<{
	username: string;
}>();

const breadcrumb = computed(() => [START, USERS, USER_DETAIL(props.username)]);
const is_own = computed(() => get_username() == props.username);
const show_collections = computed(() => is_own.value && has_permission("ADMIN"));

const search_body = (limit?: number, offset?: number): SchemaSearchRequest => ({
	user_name: props.username,
	limit: limit,
	offset: offset,
	complete: true
});

const fetch_entries = (limit: number, offset: number) => schema_search(search_body(limit, offset));
const count_entries = () => schema_count(search_body());
</script>

<template>
	<bread-crumb :items="breadcrumb"/>
	<user-profile :username="username" :key="username">
		<template #actions v-if="is_own">
			<router-link to="/import" class="btn btn-sm btn-success"><i class="fa fa-upload"></i> Import files</router-link>
			<router-link :to="'/collections/' + username" class="btn btn-sm btn-outline-light" v-if="show_collections">
				<i class="fa fa-object-group"></i> Collections
			</router-link>
			<router-link to="/profile" class="btn btn-sm btn-outline-light"><i class="fa fa-gear"></i> Account settings</router-link>
		</template>
	</user-profile>
	<h2 class="bx-section-title">{{ is_own ? 'Your schematics' : 'Schematics' }}</h2>
	<paged-content
		:key="username"
		:fetch_entries="fetch_entries"
		:count_entries="count_entries"
		:per_page="24">
		<template #body="{ list }">
			<schema-list :list="list" :show_author="false">
				<template #empty>
					<empty-state icon="cubes" title="No schematics yet" v-if="is_own">
						<p>Save a build in-game with <code>/bx_save</code> or <router-link to="/import">import a file</router-link>.</p>
						<router-link to="/mod">How to upload</router-link>
					</empty-state>
					<empty-state icon="cubes" :title="username + ' hasn\'t shared any schematics yet'" v-else/>
				</template>
			</schema-list>
		</template>
	</paged-content>
</template>
