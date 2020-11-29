import "./components/start-page.js";
import "./components/nav-bar.js";
import "./components/app.js";

const router = new VueRouter({
  routes: [{
    path: "/", component: { template: "<start-page/>" }
  }]
});

new Vue({
  el: "#app",
  router: router
});
