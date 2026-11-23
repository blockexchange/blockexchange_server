import { init, execute, is_supported, set_basedir, set_proxy, addPack, set_progress, show_error, wait_for_start } from "./wasm_helper.js";

// "Edit Online": starts the luanti wasm client with the game and mods of a schematic
// usage: /wasm/edit.html?schema_uid=<uid>&back=<path of the schematic page>

async function main() {
    const schema_uid = new URLSearchParams(location.search).get("schema_uid");
    if (!schema_uid) {
        throw new Error("no schematic selected");
    }
    if (!is_supported()) {
        throw new Error("your browser doesn't support SharedArrayBuffer, which the client needs");
    }

    const uid = encodeURIComponent(schema_uid);
    set_progress("schematic");
    const r = await fetch(`/api/schema/${uid}`);
    if (!r.ok) {
        throw new Error(`could not load the schematic (${r.status})`);
    }
    const schema = await r.json();
    if (!schema) {
        throw new Error("schematic not found");
    }
    document.title = `${schema.name} - Edit Online`;
    document.getElementById("schema-name").innerText = schema.name;
    // link back to the schematic page, local paths only
    const back = new URLSearchParams(location.search).get("back");
    if (back && back.startsWith("/") && !back.startsWith("//") && !back.includes("\\")) {
        document.getElementById("back-link").href = back;
    }

    // proxy for the http requests of the blockexchange mod to this server
    const ws_protocol = location.protocol == "https:" ? "wss:" : "ws:";
    set_proxy(`${ws_protocol}//${location.host}/api/schema/${uid}/wasm/ws`);
    set_basedir("/wasm/runtime");

    await init();
    await addPack("certs", null, "certificates");
    // the server clones the game and mods of the schematic, may take a while
    await addPack("schematic", `/api/schema/${uid}/wasm/pack`, "game and mods");

    await wait_for_start();
    // the world name is the schematic name (world.mt in the pack)
    execute(["--worldname", schema.name, "--go"]);
}

main().catch(e => show_error(e.message || String(e)));
