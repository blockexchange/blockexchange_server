import html from '../html.js';

import { get_claims } from '../../store/token.js';

export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;
    const claims = get_claims();

    if (claims && claims.user_id == schema.user_id && claims.permissions.schema.update){
      return html`
        <button class="secondary" href="#!/schema/${schema.user.name}/${schema.name}/edit">
          <i class="fa fa-edit"/> Edit
        </button>
      `;
    }
  }
};
