import unbox from './unbox.js';

export default function(we_schema){
  const cleanlua = we_schema.toString().substring(2);
  const ast = luaparse.parse(cleanlua);

  return unbox(ast.body[0].arguments[0]);
}
