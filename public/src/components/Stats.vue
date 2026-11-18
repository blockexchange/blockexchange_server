<script setup lang="ts">
import { computed } from 'vue';
import { get_stats } from "../service/info";
import format_size from "../util/format_size";
import format_count from "../util/format_count";

const stats = computed(get_stats);
</script>

<template>
	<div class="row g-3" v-if="stats && stats.schema_count >= 0">
		<div class="col-6 col-md-3">
			<div class="bx-stat">
				<div class="value">{{ format_count(stats.schema_count) }}</div>
				<div class="label">Schematics</div>
			</div>
		</div>
		<div class="col-6 col-md-3">
			<div class="bx-stat">
				<div class="value">{{ format_count(stats.user_count) }}</div>
				<div class="label">Builders</div>
			</div>
		</div>
		<div class="col-6 col-md-3">
			<div class="bx-stat">
				<div class="value">{{ format_count(stats.schemapart_count) }}</div>
				<div class="label">Stored parts</div>
			</div>
		</div>
		<div class="col-6 col-md-3">
			<div class="bx-stat">
				<div class="value">{{ format_size(stats.total_size) }}</div>
				<div class="label">Total data</div>
			</div>
		</div>
	</div>
</template>
