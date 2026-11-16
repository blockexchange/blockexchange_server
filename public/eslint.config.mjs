import js from '@eslint/js';
import globals from 'globals';

export default [
	{ ignores: ['dist/', 'node_modules/'] },
	js.configs.recommended,
	{
		files: ['js/**/*.js'],
		languageOptions: {
			ecmaVersion: 'latest',
			sourceType: 'module',
			globals: { ...globals.browser }
		}
	}
];
