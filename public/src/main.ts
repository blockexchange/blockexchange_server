import { createApp } from 'vue';
import { createRouter, createWebHashHistory } from 'vue-router';

import 'bootswatch/dist/darkly/bootstrap.min.css';
import '@fortawesome/fontawesome-free/css/all.min.css';

import App from './App.vue';
import routes from './routes';
import { fetch_info } from './service/info';
import { fetch_tags } from './service/tags';
import { check_login } from './service/login';

try {
	await Promise.all([check_login(), fetch_info(), fetch_tags()]);
	const router = createRouter({
		history: createWebHashHistory(),
		routes: routes
	});

	createApp(App)
	.use(router)
	.mount("#app");
} catch (e) {
	console.error(e);
	const el = document.getElementById("app");
	if (el) {
		el.textContent = `Error while starting the app: ${e}`;
	}
}
