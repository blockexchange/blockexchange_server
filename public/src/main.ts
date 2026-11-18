import { createApp } from 'vue';
import { createRouter, createWebHashHistory } from 'vue-router';

import 'bootswatch/dist/darkly/bootstrap.min.css';
import '@fortawesome/fontawesome-free/css/all.min.css';
import './style.css';

import App from './App.vue';
import routes from './routes';
import { fetch_info } from './service/info';
import { fetch_tags } from './service/tags';
import { check_login } from './service/login';

try {
	await Promise.all([check_login(), fetch_info(), fetch_tags()]);
	const router = createRouter({
		history: createWebHashHistory(),
		routes: routes,
		scrollBehavior: () => ({ top: 0 })
	});

	createApp(App)
	.use(router)
	.mount("#app");
} catch (e) {
	console.error(e);
	const el = document.getElementById("app");
	if (el) {
		el.innerHTML = "";
		const div = document.createElement("div");
		div.className = "container py-5 text-center";
		div.textContent = `Block exchange could not be loaded: ${e}. Please try again in a moment.`;
		el.appendChild(div);
	}
}
