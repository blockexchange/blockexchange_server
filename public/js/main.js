
import routes from './routes.js';
import Nav from './components/Nav.js';
import Breadcrumb from './components/Breadcrumb.js';

import { setup_tracker } from './util/matomo.js';
import info from './service/info.js';

// Render main route and static navbar
m.mount(document.getElementById("nav"), Nav);
m.mount(document.getElementById("breadcrumb"), Breadcrumb);
m.route(document.getElementById("app"), "/search", routes);

// fetch general infos and set up tracker if configured
info().then(info => {
  if (info.matomo.url && info.matomo.id) {
    setup_tracker(info.matomo.url, info.matomo.id);
  }
});
