import convert_we from '../../util/import/convert.js';
import parse_we from '../../util/import/parse.js';
import state from './state.js';

import { create as create_schema, complete as complete_schema } from '../../api/schema.js';
import { create as create_mods } from '../../api/schemamods.js';
import { create as create_schemapart } from '../../service/schemapart.js';
import { get_claims } from '../../store/token.js';


export function parse(){
  state.blocks = parse_we(state.we_schema);
}

export function convert(){
  state.result = convert_we(state.blocks);
  state.name = "WE-Import @ " + Date.now();
  state.description = "Imported @ " + moment().format("YYYY-MM-DD HH:mm:ss");
  state.license = "CC0";
  state.progress = 0;
}

export function upload(){
  const result = state.result;
  const total_count = result.parts.length;

  create_schema({
    name: state.name,
    description: state.description,
    max_x: state.result.max_x,
    max_y: state.result.max_y,
    max_z: state.result.max_z,
    license: state.license,
    part_length: 16
  })
  .then(function(schema){
    function worker(){
      const part = result.parts.shift();
      state.progress = (total_count - result.parts.length) / total_count * 100;
      m.redraw();

      if (!part){
        create_mods(schema.id, result.stats)
        .then(() => complete_schema(schema))
        .then(() => {
          state.result = null;
          state.blocks = null;
          state.progress = 0;
          window.location.hash = `#!/schema/${get_claims().username}/${schema.name}`;
        });
        return;
      }

      part.schema_id = schema.id;

      create_schemapart(part)
      .then(() => setTimeout(worker, 500));
    }

    worker();
  });
}
