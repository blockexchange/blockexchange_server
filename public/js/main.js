
import routes from './routes.js';
import Nav from './views/Nav.js';

m.mount(document.getElementById("nav"), Nav);
m.route(document.getElementById("app"), "/search", routes);
