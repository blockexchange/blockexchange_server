import html from '../html.js';

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

    this.links = [{
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
    }];

    this.load_data();
  }

  load_data(){
    get_by_user_and_schemaname(this.state.username, this.state.schemaname)
    .then(s => this.state.schema = s);
  }

  save(){
    const schema = this.state.schema;

    update(schema)
    .then(() => window.location.hash = `#!/schema/${schema.schemagroup.name}/${schema.name}`);
  }

  view(){
    const schema = this.state.schema;
    if (!schema){
      return;
    }

    return html`
      <div>
        <${Breadcrumb} links=${this.links}/>
        <${SchemaTitle} schema=${schema}/>
        <hr/>
        <div class="row">
          <div class="col-md-6">
            <input type="text"
              class="form-control"
              placeholder="Schema name"
              value=${schema.name}
              oninput=${e => schema.name = e.target.value}/>
          </div>
          <div class="col-md-6">
            <h3>${schema.name}</h3>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="col-md-6">
          <textarea class="form-control"
            value=${schema.description}
            style="height: 350px"
            oninput=${e => schema.description = e.target.value}
          />
        </div>
        <div class="col-md-6">
          <pre>${schema.description}</pre>
        </div>
      </div>
      <div class="row">
        <div class="col-md-6">
          <input type="text"
            class="form-control"
            placeholder="License"
            value=${schema.license}
            oninput=${e => schema.license = e.target.value}
          />
        </div>
        <div class="col-md-6">
          <${LicenseBadge} license=${schema.license}/>
        </div>
      </div>
      <hr/>
      <div class="row">
        <div class="col-md-12">
          <button class="btn btn-primary btn-block"
            onclick=${() => this.save()}>
            <i class="fa fa-save"> Save</i>
          </button>
        </div>
      </div>
    `;
  }

}
