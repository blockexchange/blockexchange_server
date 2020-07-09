
import state from './state.js';
import html from '../html.js';

import Configure from './Configure.js';
import Progress from './Progress.js';

import { parse, convert } from './actions.js';

import { hr } from '../fragments/html.js';
import Breadcrumb from '../Breadcrumb.js';

import { get_claims } from '../../store/token.js';

function readFile(file){
  const reader = new FileReader();
  reader.onload = function(fe){
    state.we_schema = fe.target.result;
    parse();
    convert();
    m.redraw();
  };
  reader.readAsBinaryString(file);
}

const links = [{
  name: "Home",
  link: "#!/"
},{
  name: "Import",
}];

export default {
  view: function(){
    if (!get_claims()){
      window.location.hash = "#!/";
      return;
    }

    return html`
      <${Breadcrumb} links=${links}/>
      <h3>Import Worldedit-Schema</h3>
      <input type="file"
        onchange=${e => readFile(e.target.files[0])}
      />
      <hr/>
      <${Progress}/>
      <hr/>
      ${state.blocks ? html`<${Configure}/>` : null }
    `;
  }
};
