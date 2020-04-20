import { get_by_user_and_schemaname } from '../../api/searchschema.js';

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
    .then(schema => console.log(schema));
  }

  view(){
    //TODO
  }

}
