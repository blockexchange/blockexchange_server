
import html from './html.js';

export default {
  view: function(vnode){
    switch (vnode.attrs.license) {
      case "CC0":
  			return html`<img src="pics/license_cc0.png"/>`;
      case "CC-BY":
        return html`<img src="pics/license_cc-by.png"/>`;
      case "CC-BY-SA":
        return html`<img src="pics/license_cc-by-sa.png"/>`;
      case "CC-BY-NC":
        return html`<img src="pics/license_cc-by-nc.png"/>`;
  		default:
        return html`<span class="badge badge-primary">${vnode.attrs.license}</span>`;
  	}
  }
};
