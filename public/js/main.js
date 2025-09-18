import App from './app.js';
import routes from './routes.js';
import { fetch_info } from './service/info.js';
import { fetch_tags } from './service/tags.js';
import { check_login } from './service/login.js';
import * as THREE from 'three';
import { decompressSync } from 'fflate';

console.log(THREE); //TODO
console.log(decompressSync); //TODO

function start(){
	// create router instance
	const router = VueRouter.createRouter({
		history: VueRouter.createWebHistory(),
		routes: routes
	});

	// start vue
	const app = Vue.createApp(App);
	app.use(router);
	app.mount("#app");
}

Promise.all([check_login(), fetch_info(), fetch_tags()])
.then(() => start());