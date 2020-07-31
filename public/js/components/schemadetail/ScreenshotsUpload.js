import html from '../html.js';

import { create } from '../../api/schemascreenshot.js';

function readFile(schemaid, file){
  create(schemaid, file.type, file.name, file);
}

export default {
  view: vnode => html`
      <input type="file"
        onchange=${e => readFile(vnode.attrs.schemaid, e.target.files[0])}
         accept="image/png, image/jpeg"
      />
    `
};
