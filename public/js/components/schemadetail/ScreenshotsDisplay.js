import html from '../html.js';

import { get_all } from '../../api/schemascreenshot.js';

const screenshot = (schema_id, s) => html`
  <img caption=${s.title} src=${`api/schema/${schema_id}/screenshot/${s.id}`}/>
`;

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
      ${this.state.screenshots.map(s => screenshot(vnode.attrs.schemaid, s))}
    `;
  }
}
