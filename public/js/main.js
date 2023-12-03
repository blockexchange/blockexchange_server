import App from './app.js';
import routes from './routes.js';

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

start();