import App from './app.js';
import routes from './routes.js';
import { fetch_info } from './service/info.js';
import { check_login } from './service/login.js';

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

Promise.all([check_login(), fetch_info()])
.then(() => start());