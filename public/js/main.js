import "./components/start-page.js";
import "./components/nav-bar.js";
import "./components/app.js";

import routes from './routes.js';

const router = new VueRouter({
  routes: routes
});

new Vue({
  el: "#app",
  router: router
});
