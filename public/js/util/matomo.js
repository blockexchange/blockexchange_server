
export function setup_tracker(url, id){
  window._paq = [
    'trackPageView',
    'enableLinkTracking',
    'setTrackerUrl', url + "matomo.php",
    'setSiteId', id
  ];

  const tag = document.createElement('script');
  tag.type='text/javascript';
  tag.async = true;
  tag.defer = true;
  tag.src = url + "piwik.js";

  document.body.appendChild(tag);
}
