
import routes from './routes.js';
import Nav from './components/Nav.js';

m.mount(document.getElementById("nav"), Nav);
m.route(document.getElementById("app"), "/search", routes);
