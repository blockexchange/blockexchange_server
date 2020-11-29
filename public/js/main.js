import './components.js';
import routes from './routes.js';

const router = new VueRouter({
  routes: routes
});

new Vue({
  el: "#app",
  router: router
});
