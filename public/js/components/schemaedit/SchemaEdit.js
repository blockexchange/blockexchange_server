
import { div, textarea, pre, hr } from '../fragments/html.js';
import { fa } from '../fragments/fa.js';
import { row, col6, col12 } from '../fragments/bootstrap.js';

import SchemaTitle from '../SchemaTitle.js';
import LicenseBadge from '../LicenseBadge.js';
import Breadcrumb from '../Breadcrumb.js';

import { get_by_user_and_schemaname } from '../../api/searchschema.js';
import { update } from '../../api/schema.js';


export default class {

  constructor(vnode) {
    this.state = {
      username: vnode.attrs.username,
      schemaname: vnode.attrs.schemaname,
      schema: null
    };

    this.load_data();
  }

  load_data(){
    get_by_user_and_schemaname(this.state.username, this.state.schemaname)
    .then(s => this.state.schema = s);
  }

  save(){
    update(this.state.schema)
    .then(() => window.location.hash = `#!/schema/${this.state.username}/${this.state.schemaname}`);
  }

  view(){
    const schema = this.state.schema;
    if (!schema){
      return;
    }

    return div([
      m(Breadcrumb, {
        links: [{
          name: "Home",
          link: "#!/"
        },{
          name: "User-schemas",
        },{
          name: this.state.username,
          link: "#!/schema/" + this.state.username
        },{
          name: this.state.schemaname,
          link: "#!/schema/" + this.state.username + "/" + this.state.schemaname
        },{
          name: "Edit",
          active: true
        }]
      }),
      m(SchemaTitle, { schema: schema }),
      hr(),
      row([
        col6([
          m("input", {
            class: "form-control",
            placeholder: "Schema name",
            value: schema.name,
            oninput: e => schema.name = e.target.value
          })
        ]),
        col6([
          m("h3", schema.name)
        ])
      ]),
      row([
        col6([
          textarea({
            value: schema.description,
            class: "form-control",
            style: "height: 350px",
            oninput: e => schema.description = e.target.value
          })
        ]),
        col6([
          pre(schema.description)
        ])
      ]),
      row([
        col6([
          m("input", {
            class: "form-control",
            placeholder: "License",
            value: schema.license,
            oninput: e => schema.license = e.target.value
          })
        ]),
        col6([
          m(LicenseBadge, { license: schema.license })
        ])
      ]),
      hr(),
      row(col12([
        m("button", {
          class: "btn btn-primary btn-block",
          onclick: () => this.save()
        },[
          fa("save"),
          " Save"
        ])
      ]))
    ]);
  }

}
