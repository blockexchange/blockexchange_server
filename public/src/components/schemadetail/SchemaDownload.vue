<script setup lang="ts">
import { computed, ref } from 'vue';
import type { Schema } from "../../api/types";
import ClipboardCopy from "../ClipboardCopy.vue";

const max_placement_tool_size = 100;

const props = defineProps<{
	schema: Schema;
	username: string;
}>();

const tab = ref<"online" | "offline">("online");

const can_use_placement_tool = computed(() => (
	props.schema.size_x <= max_placement_tool_size &&
	props.schema.size_y <= max_placement_tool_size &&
	props.schema.size_z <= max_placement_tool_size
));
</script>

<template>
	<ul class="nav nav-pills mb-3 gap-1">
		<li class="nav-item">
			<button type="button" class="nav-link" :class="{active: tab == 'online'}" @click="tab = 'online'">
				<i class="fa fa-globe"></i> In-game
			</button>
		</li>
		<li class="nav-item">
			<button type="button" class="nav-link" :class="{active: tab == 'offline'}" @click="tab = 'offline'">
				<i class="fa fa-file-zipper"></i> Download file
			</button>
		</li>
	</ul>

	<div v-if="tab == 'online'">
		<p class="small text-muted-bx">
			Needs the <router-link to="/mod">blockexchange mod</router-link> with web access enabled.
		</p>
		<ol class="d-flex flex-column gap-2">
			<li>Stand at the spot where the schematic should start and run <clipboard-copy text="/bx_pos1"/></li>
			<li>Optional: preview the area it will fill <clipboard-copy :text="'/bx_allocate ' + username + ' ' + schema.name"/></li>
			<li>Place it <clipboard-copy :text="'/bx_load ' + username + ' ' + schema.name"/></li>
		</ol>
		<div class="bx-note" :class="{'mb-3': can_use_placement_tool}">
			<i class="fa fa-triangle-exclamation"></i>
			Placing replaces any blocks already in that area.
		</div>
		<template v-if="can_use_placement_tool">
			<h3 class="h6 fw-bold mt-3">Placing it several times?</h3>
			<p class="mb-0">
				Get a placement tool and click wherever you want a copy:
				<clipboard-copy :text="'/bx_placer ' + username + ' ' + schema.name"/>
			</p>
		</template>
	</div>

	<div v-else>
		<div class="d-flex flex-wrap gap-2 mb-3">
			<a class="btn btn-success"
				:href="'/api/export_bx/' + schema.uid + '/' + schema.name + '.zip'">
				<i class="fa fa-download"></i>
				Block exchange file (.zip)
			</a>
			<a class="btn btn-outline-success"
				:class="{disabled: schema.total_parts >= 500}"
				:aria-disabled="schema.total_parts >= 500"
				:href="'/api/export_we/' + schema.uid + '/' + schema.name + '.we'">
				<i class="fa fa-download"></i>
				WorldEdit file (.we)
			</a>
		</div>
		<div class="bx-note info mb-3" v-if="schema.total_parts >= 500">
			<i class="fa fa-circle-info"></i>
			This schematic is too large for the WorldEdit format. Use the Block exchange file instead.
		</div>
		<div class="bx-note mb-3" v-else-if="schema.total_parts > 50">
			<i class="fa fa-triangle-exclamation"></i>
			This is a large schematic, so the WorldEdit file may be slow to import.
		</div>
		<p class="small text-muted-bx">To place the Block exchange file without an internet connection:</p>
		<ol class="d-flex flex-column gap-2 mb-0">
			<li>Put the .zip file into the <code>bxschems</code> folder of your world (<code>&lt;world&gt;/bxschems/</code>).</li>
			<li>Mark the starting point <clipboard-copy text="/bx_pos1"/></li>
			<li>Optional: preview the area <clipboard-copy :text="'/bx_allocate_local ' + schema.name"/></li>
			<li>Place it <clipboard-copy :text="'/bx_load_local ' + schema.name"/></li>
		</ol>
	</div>
</template>
