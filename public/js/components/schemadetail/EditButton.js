import html from '../html.js';

import { get_claims } from '../../store/token.js';
import { get_by_id } from '../../api/schemagroup.js';

export default {
  oncreate: function(vnode){
    const claims = get_claims();
    const user_id = claims.user_id;

    get_by_id(vnode.attrs.schema.schemagroup_id)
    .then(schemagroup => {
      this.schemagroup = schemagroup;
      this.can_update = schemagroup.permissions
        .filter(perm => perm.user_id == user_id)
        .some(perm => perm.update);
    });
  },

  view: function(vnode){
    const schema = vnode.attrs.schema;
    const claims = get_claims();

    if (this.can_update){
      return html`
        <a class="btn btn-secondary" href="#!/schema/${schema.schemagroup.name}/${schema.name}/edit">
          <i class="fa fa-edit"/> Edit
        </a>
      `;
    }
  }
};
