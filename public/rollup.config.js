const { nodeResolve } = require('@rollup/plugin-node-resolve');
const terser = require('@rollup/plugin-terser');

module.exports = [{
	input: 'js/main.js',
	output: {
		file :'js/bundle.js',
		format: 'iife',
		sourcemap: true,
		compact: true,
		plugins: [terser()]
	},
	plugins: [nodeResolve()]
}];
