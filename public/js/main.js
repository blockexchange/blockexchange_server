import App from './app.js';
import routes from './routes.js';
import infoStore from './store/info.js';
import loginService from './service/login.js';
import messages from './messages.js';
import { get_info } from './api/info.js';

function start(){
	// try to restore state
	loginService.restoreState();

	// create router instance
	const router = VueRouter.createRouter({
		history: VueRouter.createWebHashHistory(),
		routes: routes
	});

	const i18n = VueI18n.createI18n({
		fallbackLocale: 'en',
		messages: messages
	});

	// start vue
	const app = Vue.createApp(App);
	app.use(router);
	app.use(i18n);
	app.mount("#app");
}

get_info().then(info => {
	Object.keys(info).forEach(key => {
		infoStore[key] = info[key];
	});
	start(info);
});
