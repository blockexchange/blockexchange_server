# Luanti wasm runtime

Embedded luanti wasm client for the "Edit Online" page (`public/src/components/pages/SchemaEditOnline.vue`, glue code in `public/src/util/luanti_wasm.ts`), served under `/wasm/runtime/`.

* Luanti 5.17.0, build `347a7bd6fc7a` of [luanti-wasm](https://github.com/paradust7/luanti-wasm) (copied from the `luanti-wasm-playground` project, `html/5.17.0/347a7bd6fc7a`)
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

Check the glue code against the `custom` example and `launcher.js` of the playground for api changes (since 5.17.0 `emloop_init_fs` has to be called before the first pack is installed).
