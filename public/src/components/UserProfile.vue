<script setup lang="ts">
import { count_user_schema_stars, search_users } from "../api/user";
import { schema_count } from "../api/schema";
import type { User } from "../api/types";

import LoadingBlock from "./LoadingBlock.vue";

const props = defineProps<{
	username: string;
}>();

defineSlots<{
	actions?(): unknown;
}>();

async function fetch_user(): Promise<User & { stars: number }> {
	const l = await search_users({ name: props.username });
	const u = l[0]!;
	const stars = await count_user_schema_stars(u.uid);
	return { ...u, stars };
}

function fetch_data() {
	return {
		user: fetch_user(),
		schema_count: schema_count({ user_name: props.username })
	};
}

const account_type = (t: string) => t == "LOCAL" ? "Block exchange account" : `Linked with ${t.charAt(0) + t.slice(1).toLowerCase()}`;
</script>

<template>
	<loading-block :fetch_data="fetch_data" v-slot="{ data }">
		<div class="card mb-4">
			<div class="card-body d-flex flex-wrap align-items-center gap-3">
				<img v-if="data.user.avatar_url" :src="data.user.avatar_url" height="72" width="72" class="rounded" alt=""/>
				<div v-else class="rounded bg-secondary d-flex align-items-center justify-content-center" style="width: 72px; height: 72px;">
					<i class="fa fa-user fa-2x text-muted-bx"></i>
				</div>
				<div class="flex-grow-1">
					<h1 class="h3 mb-1">
						{{ data.user.name }}
						<span class="badge bg-warning text-dark fs-6 align-middle" v-if="data.user.role == 'ADMIN'">Admin</span>
					</h1>
					<div class="small text-muted-bx">
						<template v-if="data.user.created > 0">Joined {{ new Date(data.user.created).toLocaleDateString() }} · </template>
						{{ account_type(data.user.type) }}
					</div>
				</div>
				<div class="d-flex gap-4 text-center">
					<div>
						<div class="h4 mb-0">{{ data.schema_count }}</div>
						<div class="small text-muted-bx">{{ data.schema_count == 1 ? 'schematic' : 'schematics' }}</div>
					</div>
					<div>
						<div class="h4 mb-0"><i class="fa fa-star bx-star fs-6 align-middle"></i> {{ data.user.stars }}</div>
						<div class="small text-muted-bx">stars received</div>
					</div>
				</div>
				<div class="w-100 d-flex flex-wrap gap-2" v-if="$slots.actions">
					<slot name="actions"></slot>
				</div>
			</div>
		</div>
	</loading-block>
</template>
