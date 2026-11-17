<script setup lang="ts">
import format_time from "../util/format_time";
import { count_user_schema_stars, search_users } from "../api/user";
import { schema_count } from "../api/schema";
import type { User } from "../api/types";

import LoadingBlock from "./LoadingBlock.vue";

const props = defineProps<{
	username: string;
}>();

function fetch_data() {
	return {
		user: search_users({ name: props.username })
			.then(l => l[0]!)
			.then(u => count_user_schema_stars(u.uid).then((stars): User & { stars: number } => ({ ...u, stars }))),
		schema_count: schema_count({ user_name: props.username })
	};
}
</script>

<template>
	<loading-block :fetch_data="fetch_data" v-slot="{ data }">
		<div class="row">
			<div class="col-md-8">
				<h5>
					User profile for
					<small class="text-body-secondary">{{ data.user.name }}</small>
				</h5>
				<ul>
					<li>
						<b>Username:</b> {{ data.user.name }}
					</li>
					<li>
						<b>Schematics:</b> <span class="badge bg-secondary rouded-pill">{{ data.schema_count }}</span>
					</li>
					<li>
						<b>Stars:</b>
						<i class="fa fa-star" :style="{color: data.user.stars ? 'yellow' : ''}"></i>
						<span class="badge bg-secondary rouded-pill">{{ data.user.stars }}</span>
					</li>
					<li>
						<b>Created:</b> {{ format_time(data.user.created) }}
					</li>
					<li>
						<b>ID:</b> <span class="badge bg-success">{{ data.user.uid }}</span>
					</li>
					<li>
						<b>Type:</b> <span class="badge bg-secondary">{{ data.user.type }}</span>
					</li>
					<li>
						<b>Role:</b> <span class="badge bg-secondary">{{ data.user.role }}</span>
					</li>
				</ul>
			</div>
			<div class="col-md-4 text-end">
				<img v-if="data.user.avatar_url" :src="data.user.avatar_url" height="64" width="64"/>
			</div>
		</div>
	</loading-block>
</template>
