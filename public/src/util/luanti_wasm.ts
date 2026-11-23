// glue code for the luanti wasm client (emscripten build, embedded under /wasm/runtime/),
// based on the custom example of luanti-wasm-playground (5.17.0).
// The client needs SharedArrayBuffer (cross-origin isolated document) and can only be started once per document.

// globals of the emscripten runtime (luanti.js is loaded as a classic script)
declare global {
	function _malloc(size: number): number;
	function cwrap(name: string, ret: string | null, args: string[]): (...args: number[]) => void;
	function stringToNewUTF8(s: string): number;
	var HEAPU8: Uint8Array;
	var HEAPU32: Uint32Array;

	interface Window {
		Module?: object;
		emloop_ready?: () => void;
		emloop_fs_ready?: (persistent: number) => void;
		emloop_pack_progress?: (name: string, fraction: number) => void;
		emloop_pack_installed?: (name: string) => void;
		emloop_exited?: (status: number) => void;
	}
}

export const WASM_RUNTIME_URL = "/wasm/runtime";

export interface LuantiPack {
	name: string;
	// defaults to the pack of the runtime
	url?: string;
	// shown in the progress messages
	label: string;
}

export interface LuantiClientOptions {
	// needs the id "canvas": it is transferred to the render thread by id (OffscreenCanvas)
	canvas: HTMLCanvasElement;
	// websocket proxy for the network access of the client
	proxy_url: string;
	// installed after the base pack, in order
	packs: LuantiPack[];
	// command line arguments
	args: string[];
	on_progress: (msg: string) => void;
	// the game exited or crashed
	on_exit: (reason: string) => void;
}

export const is_supported = () => !!window.SharedArrayBuffer && window.crossOriginIsolated;

interface RuntimeFunctions {
	emloop_init_fs: (want_opfs: number) => void;
	emloop_invoke_main: (argc: number, argv: number) => void;
	irrlicht_resize: (width: number, height: number) => void;
	emsocket_init: () => void;
	emsocket_set_proxy: (proxy: number) => void;
}

let started = false;

export async function start_luanti_client(opts: LuantiClientOptions) {
	if (started) {
		throw new Error("the client can only be started once per page load");
	}
	started = true;
	if (opts.canvas.id != "canvas") {
		throw new Error("the canvas needs the id 'canvas'");
	}

	// exported functions of the runtime, set when it is ready
	let emloop_install_pack: (name: number, version: number, offset: number, len: number) => void = () => {};

	// pack name -> pending installation
	const pending_packs: Record<string, { resolve: () => void, label: string }> = {};

	window.emloop_pack_progress = (name, fraction) => {
		const p = pending_packs[name];
		if (p) {
			opts.on_progress(`${p.label} (installing ${Math.round(fraction * 100)}%)`);
		}
	};
	window.emloop_pack_installed = name => {
		const p = pending_packs[name];
		if (p) {
			delete pending_packs[name];
			p.resolve();
		}
	};
	window.emloop_exited = () => opts.on_exit("The game has exited.");

	const add_pack = async (pack: LuantiPack) => {
		opts.on_progress(`${pack.label} (downloading)`);
		const r = await fetch(pack.url || `${WASM_RUNTIME_URL}/packs/${pack.name}.pack`);
		if (!r.ok) {
			let msg = `${r.status} ${r.statusText}`;
			try {
				const body = await r.json();
				if (body?.message) {
					msg = body.message;
				}
			} catch {
				// not a json error response
			}
			throw new Error(`could not fetch the ${pack.label}: ${msg}`);
		}
		const data = new Uint8Array(await r.arrayBuffer());

		const installed = new Promise<void>(resolve => pending_packs[pack.name] = { resolve, label: pack.label });
		const offset = _malloc(data.byteLength);
		HEAPU8.set(data, offset);
		emloop_install_pack(stringToNewUTF8(pack.name), stringToNewUTF8(""), offset, data.byteLength);
		return installed;
	};

	// load the runtime
	const ready = new Promise<RuntimeFunctions>((resolve, reject) => {
		window.emloop_ready = () => {
			emloop_install_pack = cwrap("emloop_install_pack", null, ["number", "number", "number", "number"]);
			resolve({
				emloop_init_fs: cwrap("emloop_init_fs", null, ["number"]),
				emloop_invoke_main: cwrap("emloop_invoke_main", null, ["number", "number"]),
				irrlicht_resize: cwrap("irrlicht_resize", null, ["number", "number"]),
				emsocket_init: cwrap("emsocket_init", null, []),
				emsocket_set_proxy: cwrap("emsocket_set_proxy", null, ["number"])
			});
		};

		window.Module = {
			canvas: opts.canvas,
			print: (s: string) => {
				console.log("luanti", s);
				if (s.startsWith("main() exited with return value")) {
					opts.on_exit("The game has exited.");
				}
				if (s.startsWith("Unhandled exception:")) {
					opts.on_exit("The game crashed (unhandled exception), see the browser console for details.");
				}
			},
			printErr: (s: string) => console.warn("luanti", s)
		};

		const script_el = document.createElement("script");
		script_el.src = `${WASM_RUNTIME_URL}/luanti.js`;
		script_el.onerror = () => reject(new Error("could not load the client runtime"));
		document.body.appendChild(script_el);
	});

	opts.on_progress("client runtime");
	const { emloop_init_fs, emloop_invoke_main, irrlicht_resize, emsocket_init, emsocket_set_proxy } = await ready;
	// mount /luanti before the packs are installed. In memory (no OPFS): every start is a fresh world from the packs
	await new Promise<void>(resolve => {
		window.emloop_fs_ready = () => resolve();
		emloop_init_fs(0);
	});
	await add_pack({ name: "base", label: "base pack" });
	for (const pack of opts.packs) {
		await add_pack(pack);
	}

	// start
	const argv = _malloc((opts.args.length + 2) * 4);
	["./luanti", ...opts.args].forEach((arg, i) => HEAPU32[(argv >>> 2) + i] = stringToNewUTF8(arg));
	HEAPU32[(argv >>> 2) + opts.args.length + 1] = 0;
	emloop_invoke_main(opts.args.length + 1, argv);

	emsocket_init();
	emsocket_set_proxy(stringToNewUTF8(opts.proxy_url));

	const resize = () => irrlicht_resize(opts.canvas.clientWidth, opts.canvas.clientHeight);
	resize();
	window.addEventListener("resize", resize);

	block_browser_keys(opts.canvas);
}

// The key handlers of SDL run on the game thread: emscripten forwards the events asynchronously
// and never calls preventDefault(), so the browser acts on every key ("/" opens the quick find in firefox,
// F5 reloads the page, ...). Only the default action is prevented here, the events still reach SDL.
function block_browser_keys(canvas: HTMLCanvasElement) {
	// leave the overlay buttons usable with the keyboard
	const is_game_target = (e: Event) => e.target == canvas || e.target == document.body;

	// printable keys: the quick find ("/", "'") is triggered on keypress.
	// Not on keydown, that would suppress the keypress event which SDL needs for the text input
	window.addEventListener("keypress", e => {
		if (is_game_target(e)) {
			e.preventDefault();
		}
	}, true);

	// non-printable keys (F1-F12, Tab, Backspace, arrows, ...), ctrl/meta combinations keep their browser function
	window.addEventListener("keydown", e => {
		if (is_game_target(e) && e.key.length > 1 && !e.ctrlKey && !e.metaKey) {
			e.preventDefault();
		}
	}, true);
}
