
import routes from './routes.js';
import Nav from './views/Nav.js';

import { get_info } from './api/info.js';

m.mount(document.getElementById("nav"), Nav);
m.route(document.getElementById("app"), "/search", routes);

get_info()
.then(info => {
  //TODO: cleanup
  if (info.matomo.url && info.matomo.id) {
    window._paq = [
      'trackPageView',
      'enableLinkTracking',
      'setTrackerUrl', info.matomo.url + "matomo.php",
      'setSiteId', info.matomo.id
    ];

    const tag = document.createElement('script');
    tag.type='text/javascript';
    tag.async = true;
    tag.defer = true;
    tag.src = info.matomo.url + "piwik.js";

    document.body.appendChild(tag);
  }
});
