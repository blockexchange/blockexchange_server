import { fa } from '../fragments/fa.js';

import { get_claims } from '../../store/token.js';
import { remove } from '../../api/schema.js';

export default class {
  constructor(vnode){
    this.state = {
      schema: vnode.attrs.schema,
      confirm: false
    };
  }

  view(){
    const schema = this.state.schema;
    const claims = get_claims();

    if (claims && claims.user_id == schema.user_id && claims.permissions.schema.delete){
      if (!this.state.confirm){
        // delete button
        return m("button", {
          class: "btn btn-danger",
          onclick: () => this.state.confirm = true
        }, [fa("trash"), " Delete"]);
      } else {
        // confirmation buttons
        return [
          m("button", {
            class: "btn btn-danger",
            onclick: () => remove(schema)
            .then(() => window.location.hash = `#!/schema/${schema.user.name}`)
          }, "Confirm"),
          m("button", {
            class: "btn btn-success",
            onclick: () => this.state.confirm = false
          }, "Cancel"),
        ];
      }
    }
  }
}
