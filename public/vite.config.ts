import { defineConfig, type Plugin } from 'vite';
import vue from '@vitejs/plugin-vue';
import { writeFileSync, mkdirSync } from 'node:fs';

// go:embed needs a non-empty dist directory
const keepDist = (): Plugin => ({
	name: 'keep-dist',
	closeBundle() {
		mkdirSync('dist', { recursive: true });
		writeFileSync('dist/.gitkeep', '');
	}
});

// the luanti wasm launcher needs a cross-origin isolated page (SharedArrayBuffer),
// limited to /wasm/ like on the go server
const wasmHeaders = (): Plugin => ({
	name: 'wasm-headers',
	configureServer(server) {
		server.middlewares.use((req, res, next) => {
			if (req.url?.startsWith('/wasm/')) {
				res.setHeader('Cross-Origin-Opener-Policy', 'same-origin');
				res.setHeader('Cross-Origin-Embedder-Policy', 'require-corp');
				res.setHeader('Cross-Origin-Resource-Policy', 'same-origin');
			}
			next();
		});
	}
});

export default defineConfig({
	base: '/',
	publicDir: 'static',
	plugins: [vue(), keepDist(), wasmHeaders()],
	build: {
		outDir: 'dist',
		emptyOutDir: true,
		sourcemap: true
	},
	server: {
		proxy: {
			// websocket proxy of the wasm client, keeps the host header for the same-origin check
			'^/api/schema/[^/]+/wasm/ws': {
				target: process.env.API_TARGET || 'http://localhost:8080',
				changeOrigin: false,
				ws: true
			},
			'/api': {
				target: process.env.API_TARGET || 'http://localhost:8080',
				changeOrigin: true
			},
			// luanti wasm build, embedded in the go server
			'/wasm/runtime': {
				target: process.env.API_TARGET || 'http://localhost:8080',
				changeOrigin: true
			}
		}
	}
});
