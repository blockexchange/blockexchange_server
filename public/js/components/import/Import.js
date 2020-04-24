
import { create as create_schema, complete as complete_schema } from '../../api/schema.js';
import { create as create_schemapart } from '../../api/schemapart.js';

import unbox from './unbox.js';
import convert from './convert.js';



function upload_parts(result, name, description){

  create_schema({
    name: name,
    description: description,
    size_x: result.max_x,
    size_y: result.max_y,
    size_z: result.max_z,
    license: "CC0",
    part_length: 16
  })
  .then(function(schema){
    function worker(){
      const part = result.parts.shift();
      if (!part){
        complete_schema(schema);
        return;
      }

      part.schema_id = schema.id;

      console.log(part);
      create_schemapart(part)
      .then(() => setTimeout(worker, 500));
    }

    worker();
  });

}

function readString(we_schema){
  const cleanlua = we_schema.toString().substring(2);
  const ast = luaparse.parse(cleanlua);

  const blocks = unbox(ast.body[0].arguments[0]);
  console.log("unboxed", blocks.length);

  const result = convert(blocks);
  console.log("total schema parts: " + result.parts.length);
  console.log("max sizes: " + result.max_x + "/" + result.max_y + "/" + result.max_z);

  upload_parts(result, "WE-Import @" + Date.now(), "WE Import wip");
}

function readFile(file){
  const reader = new FileReader();
  reader.onload = function(fe){
    readString(fe.target.result);
  };
  reader.readAsBinaryString(file);
}

export default {
  view: function(){
    return m("input[type=file]", {
      onchange: function(e){
        const fileList = e.target.files;
        readFile(fileList[0]);
      }
    });
  }
};
