import { createApp } from 'vue';
import { createRouter, createWebHistory } from 'vue-router';

import 'bootswatch/dist/darkly/bootstrap.min.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

import App from './app.js';
import routes from './routes.js';
import { fetch_info } from './service/info.js';
import { fetch_tags } from './service/tags.js';
import { check_login } from './service/login.js';

function start(){
	// create router instance
	const router = createRouter({
		history: createWebHistory(),
		routes: routes
	});

	// start vue
	const app = createApp(App);
	app.use(router);
	app.mount("#app");
}

Promise.all([check_login(), fetch_info(), fetch_tags()])
.then(() => start());
