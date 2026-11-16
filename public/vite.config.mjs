import { defineConfig } from 'vite';
import { writeFileSync, mkdirSync } from 'node:fs';

// go:embed needs a non-empty dist directory
const keepDist = () => ({
	name: 'keep-dist',
	closeBundle() {
		mkdirSync('dist', { recursive: true });
		writeFileSync('dist/.gitkeep', '');
	}
});

export default defineConfig({
	base: '/',
	publicDir: 'static',
	plugins: [keepDist()],
	resolve: {
		alias: {
			// components use runtime-compiled string templates
			vue: 'vue/dist/vue.esm-bundler.js'
		}
	},
	build: {
		outDir: 'dist',
		emptyOutDir: true,
		sourcemap: true
	},
	server: {
		proxy: {
			'/api': {
				target: process.env.API_TARGET || 'http://localhost:8080',
				changeOrigin: true
			}
		}
	}
});
