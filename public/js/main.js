import './app.js';
import routes from './routes.js';
import infoStore from './store/info.js';
import loginService from './service/login.js';
import messages from './messages.js';

import './util/prettysize-filter.js';

function start(){
	// try to restore state
	loginService.restoreState();

	// create router instance
	const router = new VueRouter({
	  routes: routes
	});

	const i18n = new VueI18n({
	  messages: messages
	});

	// start vue
	new Vue({
	  el: "#app",
	  router: router,
		i18n: i18n
	});
}

fetch("api/info")
.then(r => r.json())
.then(info => {
	Object.keys(info).forEach(key => {
		infoStore[key] = info[key];
	});
	start();
});
