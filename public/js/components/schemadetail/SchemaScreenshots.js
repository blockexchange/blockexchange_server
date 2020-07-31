import html from '../html.js';

import { get_all, create } from '../../api/schemascreenshot.js';

function readFile(schemaid, file){
  create(schemaid, file.type, file.name, file);
}

export default class {
  constructor(vnode) {
    this.state = {
      screenshots: null
    };

    get_all(vnode.attrs.schemaid)
    .then(sc => this.state.screenshots = sc);
  }

  view(vnode) {
    if (!this.state.screenshots){
      return html`<div>Loading screenshots</div>`;
    }

    return html`
      <input type="file"
        onchange=${e => readFile(vnode.attrs.schemaid, e.target.files[0])}
         accept="image/png, image/jpeg"
      />
    `;
  }
}
