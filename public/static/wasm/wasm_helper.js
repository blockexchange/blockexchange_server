/* globals _malloc, cwrap, HEAPU32, HEAPU8, stringToNewUTF8 */

// helper for the luanti wasm client, based on the custom example of luanti-wasm-playground (5.16.1)

function makeArgv(args) {
    // Assuming 4-byte pointers
    const argv = _malloc((args.length + 1) * 4);
    let i;
    for (i = 0; i < args.length; i++) {
        HEAPU32[(argv >>> 2) + i] = stringToNewUTF8(args[i]);
    }
    HEAPU32[(argv >>> 2) + i] = 0; // argv[argc] == NULL
    return [i, argv];
}

let emloop_invoke_main, emloop_install_pack;
let irrlicht_resize;
let emsocket_init, emsocket_set_proxy;

let basedir = "wasm";
let proxy = "ws://proxy:8080";

// pack name -> { resolve, label }
const pending_packs = {};

window.emloop_pack_progress = function(name, fraction) {
    const p = pending_packs[name];
    if (p) {
        set_progress(`${p.label} (installing ${Math.round(fraction * 100)}%)`);
    }
};

window.emloop_pack_installed = function(name) {
    const p = pending_packs[name];
    if (p) {
        delete pending_packs[name];
        p.resolve();
    }
};

// newer builds only
window.emloop_exited = function(status) {
    console.log("emloop_exited", {status});
    show_exited();
};

export function is_supported() {
    return !!window.SharedArrayBuffer;
}

export function set_basedir(d) {
    basedir = d;
}

export function set_proxy(p) {
    proxy = p;
}

export function set_progress(msg) {
    document.getElementById("loading-msg").innerText = msg;
}

// fetches a pack (default: from the runtime dir) and installs it into the filesystem of the client
export async function addPack(name, url, label) {
    url = url || `${basedir}/packs/${name}.pack`;
    label = label || name;

    set_progress(`${label} (downloading)`);
    const r = await fetch(url);
    if (!r.ok) {
        let msg = `${r.status} ${r.statusText}`;
        try {
            const body = await r.json();
            if (body && body.message) {
                msg = body.message;
            }
        } catch {
            // not a json error response
        }
        throw new Error(`could not fetch ${label}: ${msg}`);
    }
    const u8arr = new Uint8Array(await r.arrayBuffer());

    const installed = new Promise(resolve => pending_packs[name] = { resolve, label });
    const len = u8arr.byteLength;
    const offset = _malloc(len);
    HEAPU8.set(u8arr, offset);
    emloop_install_pack(stringToNewUTF8(name), stringToNewUTF8(""), offset, len);
    return installed;
}

const ready = new Promise((resolve, reject) => {
    window.emloop_ready = function() {
        emloop_invoke_main = cwrap("emloop_invoke_main", null, ["number", "number"]);
        emloop_install_pack = cwrap("emloop_install_pack", null, ["number", "number", "number"]);
        irrlicht_resize = cwrap("irrlicht_resize", null, ["number", "number"]);
        emsocket_init = cwrap("emsocket_init", null, []);
        emsocket_set_proxy = cwrap("emsocket_set_proxy", null, ["number"]);

        addPack("base", null, "base pack")
        .then(() => resolve())
        .catch(reject);
    };
});

// loads the client runtime and the base pack
export function init() {
    // provide global callbacks
    window.Module = {
        canvas: document.getElementById("canvas"),
        print: s => {
            console.log("print", s);
            if (s.startsWith("main() exited with return value")) {
                show_exited();
            }
            if (s.startsWith("Unhandled exception:")) {
                show_error("unhandled exception, see the browser console for details");
            }
        },
        printErr: s => {
            console.warn("printErr", s);
        }
    };

    set_progress("client runtime");
    const script_el = document.createElement("script");
    script_el.src = `${basedir}/luanti.js`;
    script_el.onerror = () => show_error("could not load the client runtime");
    document.body.appendChild(script_el);

    return ready;
}

// shows the start button and waits for the click (fullscreen needs a user gesture)
export function wait_for_start() {
    document.getElementById("loading").style.display = "none";
    const start_el = document.getElementById("start");
    start_el.style.display = "block";

    return new Promise(resolve => {
        document.getElementById("start-btn").addEventListener("click", () => {
            start_el.remove();
            const el = document.documentElement;
            if (el.requestFullscreen) {
                // the canvas fills the window anyway if fullscreen is refused
                el.requestFullscreen().catch(e => console.warn("fullscreen refused", e));
            }
            resolve();
        }, { once: true });
    });
}

function resize() {
    irrlicht_resize(window.innerWidth, window.innerHeight);
}

export function execute(args) {
    const loading_el = document.getElementById("loading");
    if (loading_el) {
        loading_el.remove();
    }

    const canvas_el = document.getElementById("canvas");
    canvas_el.style.display = "block";

    const [argc, argv] = makeArgv(["./luanti", ...args]);
    emloop_invoke_main(argc, argv);

    emsocket_init();
    emsocket_set_proxy(stringToNewUTF8(proxy));

    resize();
    window.addEventListener("resize", resize);
}

export function show_error(err) {
    console.error(err);
    for (const id of ["canvas", "loading", "start"]) {
        const el = document.getElementById(id);
        if (el) {
            el.remove();
        }
    }
    if (document.fullscreenElement) {
        document.exitFullscreen();
    }
    document.getElementById("error").style.display = "block";
    document.getElementById("error-msg").innerText = err;
}

function show_exited() {
    show_error("The game has exited.");
}
