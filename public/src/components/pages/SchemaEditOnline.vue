<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { schema_search } from "../../api/schema";
import { get_user_uid, has_permission } from "../../service/login";
import { is_supported, start_luanti_client } from "../../util/luanti_wasm";

// "Edit Online": the luanti wasm client with the game and mods of the schematic, filling the whole page.
// Only reachable in the cross-origin isolated document under /wasm/ (see the router guard in main.ts)

const props = defineProps<{
	username: string;
	name: string;
}>();

const canvas = ref<HTMLCanvasElement>();
const progress = ref("schematic");
const error = ref("");
const running = ref(false);

// full page load, tears down the client and leaves the isolated document
const detail_url = computed(() => `/#/schema/${encodeURIComponent(props.username)}/${encodeURIComponent(props.name)}`);
const close = () => location.assign(detail_url.value);
const restart = () => location.reload();

async function start() {
	if (!is_supported()) {
		throw new Error("your browser doesn't support SharedArrayBuffer, which the client needs");
	}

	const r = await schema_search({ schema_name: props.name, user_name: props.username });
	const schema = r[0]?.schema;
	if (!schema) {
		throw new Error(`the schematic ${props.name} by ${props.username} doesn't exist`);
	}
	if (schema.user_uid != get_user_uid() && !has_permission("ADMIN")) {
		throw new Error("only the owner of the schematic can edit it online");
	}
	document.title = `${schema.name} - Edit Online`;

	const uid = encodeURIComponent(schema.uid);
	const ws_protocol = location.protocol == "https:" ? "wss:" : "ws:";
	await start_luanti_client({
		canvas: canvas.value!,
		// http requests of the blockexchange mod to this server
		proxy_url: `${ws_protocol}//${location.host}/api/schema/${uid}/wasm/ws`,
		packs: [
			{ name: "certs", label: "certificates" },
			// the server clones the game and mods of the schematic on the first start, cached afterwards
			{ name: "schematic", url: `/api/schema/${uid}/wasm/pack`, label: "game and mods" },
			// the blockexchange mod and its dependencies (worldmods), without the mods of the schematic
			{ name: "builtin", url: `/api/schema/${uid}/wasm/builtin`, label: "blockexchange mod" },
			// world.mt, map_meta.txt and minetest.conf, matching the packs above (installed last)
			{ name: "config", url: `/api/schema/${uid}/wasm/config`, label: "world settings" }
		],
		// the world folder of the packs, world.mt has no world_name so the folder name is used
		args: ["--worldname", "world", "--config", "/luanti/minetest.conf", "--go"],
		on_progress: msg => progress.value = msg,
		on_exit: reason => error.value = reason
	});
	running.value = true;
	canvas.value!.focus();
}

onMounted(() => {
	document.body.style.overflow = "hidden";
	start().catch(e => error.value = (e as Error).message || String(e));
});

onBeforeUnmount(() => {
	document.body.style.overflow = "";
});
</script>

<template>
	<div class="bx-wasm">
		<!-- the canvas must not have a border or padding, or mouse coords will be wrong.
			The id is required: the runtime hands the canvas to the render thread (OffscreenCanvas) by id
			and looks it up as "#canvas" -->
		<canvas ref="canvas" id="canvas" class="bx-wasm-canvas" tabindex="-1" @contextmenu.prevent></canvas>

		<div class="bx-wasm-overlay" v-if="error">
			<div class="alert alert-warning mb-0" style="max-width: 600px;">
				<i class="fa fa-triangle-exclamation"></i> {{ error }}
			</div>
			<div class="d-flex gap-2">
				<button type="button" class="btn btn-outline-light" @click="restart">
					<i class="fa fa-rotate-right"></i> Restart
				</button>
				<button type="button" class="btn btn-light" @click="close">
					<i class="fa fa-arrow-left"></i> Back to the schematic
				</button>
			</div>
		</div>

		<div class="bx-wasm-overlay" v-else-if="!running">
			<div class="spinner-border" role="status"></div>
			<div>Loading <b>{{ name }}</b>: {{ progress }}</div>
			<div class="small opacity-75">The game and mods are prepared on the server, the first start can take a minute.</div>
			<button type="button" class="btn btn-sm btn-outline-light" @click="close">Cancel</button>
		</div>

		<button type="button" class="btn btn-sm btn-dark bx-wasm-close" v-else @click="close" title="Close and go back to the schematic">
			<i class="fa fa-xmark"></i>
			<span class="visually-hidden">Close</span>
		</button>
	</div>
</template>

<style scoped>
.bx-wasm {
	position: fixed;
	inset: 0;
	z-index: 2000;
	background: #000;
	color: #ddd;
}
.bx-wasm-canvas {
	display: block;
	width: 100%;
	height: 100%;
	border: 0;
	padding: 0;
	outline: none;
}
.bx-wasm-overlay {
	position: absolute;
	inset: 0;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: 16px;
	padding: 16px;
	text-align: center;
	background: #111;
}
.bx-wasm-close {
	position: absolute;
	top: 8px;
	right: 8px;
	opacity: 0.6;
}
.bx-wasm-close:hover {
	opacity: 1;
}
</style>
