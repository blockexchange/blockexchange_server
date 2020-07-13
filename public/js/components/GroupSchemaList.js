import SchemaList from './schemalist/SchemaList.js';
import Breadcrumb from './Breadcrumb.js';

import html from './html.js';

import { find_by_username } from '../api/searchschema.js';

export default class {
  constructor(vnode) {
    this.state = {
      list: [],
      username: vnode.attrs.username
    };

    this.state.links = [{
      name: "Home",
      link: "#!/"
    },{
      name: "Group",
    },{
      name: this.state.username,
      active: true
    }];

    this.search();
  }

  search(){
    find_by_username(this.state.username)
    .then(l => this.state.list = l);
  }

  view(){
    return html`
      <${Breadcrumb} links=${this.state.links}/>
      <${SchemaList} list=${this.state.list} hide_group=${true}/>
    `;
  }
}
