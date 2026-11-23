import js from '@eslint/js';
import globals from 'globals';
import pluginVue from 'eslint-plugin-vue';
import tseslint from 'typescript-eslint';
import { defineConfig } from 'eslint/config';

export default defineConfig(
	{ ignores: ['dist/', 'node_modules/'] },
	js.configs.recommended,
	...tseslint.configs.recommended,
	...pluginVue.configs['flat/essential'],
	{
		files: ['**/*.vue'],
		languageOptions: {
			parserOptions: { parser: tseslint.parser }
		}
	},
	{
		files: ['src/**/*.{ts,vue}'],
		languageOptions: { globals: { ...globals.browser } },
		rules: {
			// page components are single words (Login, Users, ...)
			'vue/multi-word-component-names': 'off'
		}
	},
	{
		// plain browser scripts of the luanti wasm launcher (served as-is)
		files: ['static/wasm/**/*.js'],
		languageOptions: { globals: { ...globals.browser } }
	},
	{
		files: ['*.config.{ts,mjs}'],
		languageOptions: { globals: { ...globals.node } }
	}
);
