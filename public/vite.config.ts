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

export default defineConfig({
	base: '/',
	publicDir: 'static',
	plugins: [vue(), keepDist()],
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
