import html from '../html.js';

import { get_by_user_and_schemaname } from '../../api/searchschema.js';

import Scene from './Scene.js';

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

  view() {
    if (!this.state.schema){
      return html`<div>Loading...</div>`;
    }

    return html`
      <${Scene} schema=${this.state.schema}
        progressCallback=${f => this.state.progress = f * 100}/>
      <div class="progress">
        <div class="progress-bar" style="width: ${this.state.progress}%">
          ${Math.floor(this.state.progress * 10) / 10}%
        </div>
      </div>`;
  }
}
