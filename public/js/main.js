import './app.js';
import routes from './routes.js';
import infoStore from './store/info.js';

function start(){
	const router = new VueRouter({
	  routes: routes
	});

	new Vue({
	  el: "#app",
	  router: router
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
