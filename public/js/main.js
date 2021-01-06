import './app.js';
import routes from './routes.js';
import infoStore from './store/info.js';
import loginService from './service/login.js';
import messages from './messages.js';
import { get_info } from './api/info.js';

import './util/prettysize-filter.js';

function start(info){
	// try to restore state
	loginService.restoreState();

	// create router instance
	const router = new VueRouter({
	  routes: routes
	});

	if (info.matomo.url){
		// Enable analytics
		Vue.use(VueMatomo, {
			host: info.matomo.url,
			siteId: info.matomo.id,
			router: router
		});
	}

	const i18n = new VueI18n({
		fallbackLocale: 'en',
		messages: messages
	});

	// start vue
	new Vue({
	  el: "#app",
	  router: router,
		i18n: i18n
	});
}

get_info().then(info => {
	Object.keys(info).forEach(key => {
		infoStore[key] = info[key];
	});
	start(info);
});
