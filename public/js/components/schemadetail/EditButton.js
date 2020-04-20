import { button } from '../fragments/bootstrap.js';
import { fa } from '../fragments/fa.js';

import { get_claims } from '../../store/token.js';

export default {
  view: function(vnode){
    const schema = vnode.attrs.schema;
    const claims = get_claims();

    if (claims && claims.user_id == schema.user_id && claims.permissions.schema.update){
      return button("secondary",
        `#!/schema/${schema.user.name}/${schema.name}/edit`, [
        fa("edit"),
        " Edit"
      ]);
    }
  }
};
