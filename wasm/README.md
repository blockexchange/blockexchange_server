# Luanti wasm runtime

Embedded luanti wasm client for the "Edit Online" page (`public/static/wasm/edit.html`), served under `/wasm/runtime/`.

* Luanti 5.16.1, build `26247b229671` of [luanti-wasm](https://github.com/paradust7/luanti-wasm) (copied from the `luanti-wasm-playground` project, `html/5.16.1/26247b229671`)
* `luanti.js`, `luanti.wasm`: stored gzip-compressed (`*.gz`), served with `Content-Encoding: gzip` or decompressed for clients without gzip support
* `worker.js`: pthread worker
* `packs/base.pack`, `packs/certs.pack`: base filesystem and CA certificates (zstd compressed tar)

Luanti is licensed under the LGPL-2.1-or-later, see https://github.com/luanti-org/luanti

## Updating

```sh
R=<release dir of a luanti-wasm build>
cp $R/worker.js wasm/runtime/
gzip -9 -n -c $R/luanti.js > wasm/runtime/luanti.js.gz
gzip -9 -n -c $R/luanti.wasm > wasm/runtime/luanti.wasm.gz
cp $R/packs/base.pack $R/packs/certs.pack wasm/runtime/packs/
```
