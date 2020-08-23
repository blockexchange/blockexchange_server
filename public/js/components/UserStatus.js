import { get_claims } from '../store/token.js';

import html from './html.js';

export default {
  view: function() {
    const claims = get_claims();
    if (!claims){
      return html`<div/>`;
    }

    return html`<span class="badge badge-light">${claims.username}</span>"`;
  }
};
